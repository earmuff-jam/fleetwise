package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	stormRider "github.com/earmuff-jam/ciri-stormrider"
	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/earmuff-jam/fleetwise/service"
)

// Signup ...
// swagger:route POST /api/v1/signup Authentication signup
//
// # Sign up users into the database system.
//
// Parameters:
//   - +name: email
//     in: query
//     description: The email address of the current user
//     type: string
//     required: true
//   - +name: password
//     in: query
//     description: The password of the current user
//     type: string
//     required: true
//   - +name: birthday
//     in: query
//     description: The birthdate of the current user. Must be 13 years of age.
//     type: string
//     required: true
//   - +name: role
//     in: query
//     description: The user role for the application.
//     type: string
//     default: false
//     required: false
//
// Responses:
// 200: UserCredentials
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Signup(rw http.ResponseWriter, r *http.Request) {

	draftUser := &model.UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		config.Log("Unable to decode request parameters", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if len(draftUser.Email) <= 0 || len(draftUser.EncryptedPassword) <= 0 {
		config.Log("unable to decode user", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	if len(draftUser.Username) <= 3 {
		config.Log("user name is required and must be at least 4 characters in length", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	if len(draftUser.Role) <= 0 {
		draftUser.Role = "USER"
	}

	t, err := time.Parse("2006-01-02", draftUser.Birthday)
	if err != nil {
		config.Log("Error parsing birthdate", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	// Check if the user is at least 13 years old
	age := time.Now().Year() - +t.Year()
	if age <= 13 {
		config.Log("unable to sign up user. verification failed", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
		return
	}

	backendClientUsr := os.Getenv("CLIENT_USER")
	if len(backendClientUsr) == 0 {
		config.Log("unable to retrieve user from env", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("error")
	}

	resp, err := service.RegisterUser(backendClientUsr, draftUser)
	if err != nil {
		config.Log("Unable to create new user", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// Signin ...
// swagger:route POST /api/v1/signin Authentication signin
//
// # Sign in users into the database system.
//
// Parameters:
//   - +name: email
//     in: query
//     description: The email address of the current user
//     type: string
//     required: true
//   - +name: password
//     in: query
//     description: The password of the current user
//     type: string
//     required: true
//
// Responses:
// 200: UserCredentials
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Signin(rw http.ResponseWriter, r *http.Request) {

	draftUser := &model.UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		config.Log("Unable to decode request parameters", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}
	if len(draftUser.Email) <= 0 || len(draftUser.EncryptedPassword) <= 0 {
		config.Log("Unable to decode user", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	draftUser.UserAgent = r.UserAgent()

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		config.Log("unable to retrieve user from env. Unable to sign in.", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}

	resp, err := service.FetchUser(user, draftUser)
	if err != nil {
		config.Log("Unable to sign user in", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   draftUser.PreBuiltToken,
		Expires: draftUser.ExpirationTime,
	})

	rw.Header().Add("Role2", draftUser.Role)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// IsValidUserEmail ...
// swagger:route POST /api/v1/isValidEmail Authentication IsValidUserEmail
//
// # Returns true or false is the user is already in the system
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func IsValidUserEmail(rw http.ResponseWriter, r *http.Request) {

	draftUserEmail := &model.UserResponse{}
	err := json.NewDecoder(r.Body).Decode(draftUserEmail)
	r.Body.Close()
	if err != nil {
		config.Log("unable to validate user email address", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		config.Log("unable to retrieve user from env. Unable to sign in.", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}

	resp, err := db.IsValidUserEmail(user, draftUserEmail.EmailAddress)
	if err != nil {
		config.Log("unable to verify user email address", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// ResetPassword ...
// swagger:route POST /api/v1/resetPassword Authentication ResetPassword
//
// # Sends an email notification with email address in the jwt and token in the email address.
// Use this token to reset password.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func ResetPassword(rw http.ResponseWriter, r *http.Request) {

	draftUser := &model.UserResponse{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		config.Log("unable to reset password", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	if len(draftUser.EmailAddress) <= 0 {
		config.Log("unable to reset password", errors.New("missing required fields"))
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		config.Log("unable to retrieve user from env. Unable to sign in.", nil)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode("unable to retrieve user from env")
		return
	}

	err = service.ResetPassword(user, draftUser)
	if err != nil {
		config.Log("unable to reset password", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("200 OK")
}

// VerifyEmailAddress ...
// swagger:route GET /api/v1/verify Authentication VerifyEmailAddress
//
// # Used to verify if the user correctly verified the selected email address. If the token
// is valid, then the user was successfully verified. Updates the db with verification check.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func VerifyEmailAddress(rw http.ResponseWriter, r *http.Request) {

	user := os.Getenv("CLIENT_USER")
	if len(user) == 0 {
		config.Log("unable to retrieve client user", errors.New(config.ErrorFetchingCurrentUser))
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorFetchingCurrentUser)
		return
	}

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
		secretToken = ""
	}

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		config.Log("unable to validate request params", errors.New(config.ErrorTokenValidation))
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorTokenValidation)
		return
	}

	isValid, err := stormRider.ValidateJWT(tokenString, secretToken)

	if err != nil || !isValid {
		config.Log("unable to validate token", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorTokenValidation)
		return
	}

	resp, err := stormRider.ParseJwtToken(tokenString, secretToken)
	if err != nil {
		config.Log("unable to validate token", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorTokenValidation)
		return
	}

	draftUserID := resp.Claims.Subject
	if len(draftUserID) <= 0 {
		config.Log("unable to validate token", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorTokenSubjectValidation)
		return
	}

	err = db.VerifyUser(user, draftUserID)
	if err != nil {
		config.Log("unable to complete verification of the user", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("Verified. Return to application to sign in.")

}

// ResetEmailToken ...
// swagger:route POST /api/v1/reset Authentication ResetEmailToken
//
// # Resets the token in the database and allows users to resend email in case the token
// is incorrect or failed to reach the user. This route is kept under the secure route
// because a user must be logged in before they can activate a verification token. This is
// used to validate user email address.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func ResetEmailToken(rw http.ResponseWriter, r *http.Request, user string) {

	draftUser := &model.UserResponse{}
	err := json.NewDecoder(r.Body).Decode(draftUser)
	r.Body.Close()
	if err != nil {
		config.Log("unable to validate user details", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	if len(draftUser.EmailAddress) <= 0 {
		config.Log("unable to validate user email address", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	if len(draftUser.ID) <= 0 {
		config.Log("unable to validate user id", err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(err)
		return
	}

	if draftUser.IsVerified {
		config.Log("duplicate request detected", errors.New(config.ErrorUserIsAlreadyVerified))
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(config.ErrorUserIsAlreadyVerified)
		return
	}

	// draft user id is required so that the jwt token can be associated with the user
	service.PerformEmailNotificationService(user, draftUser.EmailAddress, draftUser.ID, config.EmailVerificationTokenStringURI)

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("200 OK")
}

// Logout ...
// swagger:route POST /api/v1/logout Authentication logout
//
// # Logs users out of the database system.
//
// Responses:
// 200: MessageResponse
// 400: MessageResponse
// 404: MessageResponse
// 500: MessageResponse
func Logout(rw http.ResponseWriter, r *http.Request) {

	// immediately clear the token cookie
	http.SetCookie(rw, &http.Cookie{
		Name:     "token",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(nil)
}
