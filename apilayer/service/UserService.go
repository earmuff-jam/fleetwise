package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	stormRider "github.com/earmuff-jam/ciri-stormrider"

	ether "github.com/earmuff-jam/ether"
	etherTypes "github.com/earmuff-jam/ether/types"

	"github.com/earmuff-jam/ciri-stormrider/types"
	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	"github.com/earmuff-jam/fleetwise/utils"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// FetchUser ...
//
// Function is used to retrieve user details and perform jwt maniupulation in the application
func FetchUser(user string, draftUser *model.UserCredentials) (*model.UserResponse, error) {

	draftTime := os.Getenv("TOKEN_VALIDITY_TIME")
	if len(draftTime) <= 0 {
		config.Log("unable to find token validity time. defaulting to default values", nil)
		draftTime = config.DefaultTokenValidityTime
	}

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
		secretToken = ""
	}

	draftUser, err := db.RetrieveUser(user, draftUser)
	if err != nil {
		config.Log("unable to retrieve user details", err)
		return nil, err
	}

	formattedTime, err := strconv.ParseInt(draftTime, 10, 64)
	if err != nil {
		config.Log("unable to parse provided time", err)
		return nil, err
	}

	draftCredentials := types.Credentials{
		Claims: jwt.StandardClaims{
			ExpiresAt: formattedTime,
			Subject:   draftUser.ID.String(),
		},
	}

	userCredsWithToken, err := stormRider.CreateJWT(&draftCredentials, secretToken)

	if err != nil {
		config.Log("unable to create JWT token", err)
		return nil, err
	}
	draftUser.PreBuiltToken = userCredsWithToken.Cookie
	draftUser.LicenceKey = userCredsWithToken.LicenceKey

	err = updateJwtToken(user, draftUser)
	if err != nil {
		config.Log("unable to upsert token", err)
		return nil, err
	}

	return &model.UserResponse{
		ID:           draftUser.ID.String(),
		EmailAddress: draftUser.Email,
		IsVerified:   draftUser.IsVerified,
	}, nil
}

// updateJwtToken ...
//
// allows to update the token schema with the proper credentials for the user
// also updates the auth.users table with the used license key to decode the jwt.
// the result however is a masked entity to preserve the users jwt
func updateJwtToken(user string, draftUser *model.UserCredentials) error {
	db, err := db.SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to setup transaction for db", err)
		return err
	}

	err = upsertLicenseKey(draftUser.ID.String(), draftUser.LicenceKey, tx)
	if err != nil {
		config.Log("unable to add license key", err)
		tx.Rollback()
		return err
	}

	err = upsertOauthToken(draftUser, tx)
	if err != nil {
		config.Log("unable to add auth token", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		config.Log("unable to commit selected transaction", err)
		return err
	}

	return nil
}

// upsertLicenseKey...
//
// set the instance_id as the license key used to encode / decode the jwt.
// save each users own key so that we can decode the token in private if needed be.
func upsertLicenseKey(userID string, licenseKey string, tx *sql.Tx) error {

	sqlStr := "UPDATE auth.users SET instance_id = $1 WHERE id = $2;"
	_, err := tx.Exec(sqlStr, licenseKey, userID)
	if err != nil {
		config.Log("unable to add license key to signed in user", err)
		tx.Rollback()
		return err
	}
	return nil
}

// upsertOauthToken ...
//
// updates the oauth token table in the database
func upsertOauthToken(draftUser *model.UserCredentials, tx *sql.Tx) error {

	var maskedID string

	sqlStr := `
	INSERT INTO auth.oauth
	(token, user_id, expiration_time, user_agent)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id)
	DO UPDATE SET
		token = EXCLUDED.token,
		expiration_time = EXCLUDED.expiration_time,
		user_agent = EXCLUDED.user_agent
	RETURNING id;`

	err := tx.QueryRow(
		sqlStr,
		draftUser.PreBuiltToken,
		draftUser.ID,
		draftUser.ExpirationTime,
		draftUser.UserAgent,
	).Scan(&maskedID)

	if err != nil {
		config.Log("unable to add token", err)
		tx.Rollback()
		return err
	}

	// apply the masked token
	draftUser.PreBuiltToken = maskedID
	return nil
}

// RegisterUser ...
//
// Performs saveUser operation and sends email service to the user to verify registration.
// Attempts to save the user in the db and sends a validity token to the user email address
// for verification purposes.
func RegisterUser(user string, draftUser *model.UserCredentials) (*model.UserCredentials, error) {

	resp, err := db.SaveUser(user, draftUser)
	if err != nil {
		config.Log("unable to save user", err)
		return nil, err
	}

	PerformEmailNotificationService(user, draftUser.Email, draftUser.ID.String(), config.EmailVerificationTokenStringURI)

	return resp, nil
}

// ResetPassword ...
//
// Invokes the ability to reset password. If the email address is valid, sends email service to the
// user with a unique token.
func ResetPassword(user string, draftUser *model.UserResponse) error {

	// returns false if user is found
	resp, err := db.IsValidUserEmail(user, draftUser.EmailAddress)
	if err != nil {
		config.Log("unable to reset password", err)
		return err
	}

	if resp {
		config.Log("unable to find selected user", errors.New("user not found"))
		return errors.New("user not found")
	}

	userID, err := db.RetrieveUserDetailsByEmailAddress(user, draftUser.EmailAddress)
	if err != nil {
		config.Log("unable to find selected user id.", err)
		return errors.New("user not found")
	}

	PerformEmailNotificationService(user, draftUser.EmailAddress, userID, config.ResetPasswordTokenStringURI)
	return nil
}

// PerformEmailNotificationService ...
//
// Updates user fields in db with new token and sends email notification for email verification
// to client using Send Grid api. This function is also re-used when users attempt to re-verify the
// token if Send Grid fails to send the api. The userID is passed in the subject in the token to verify
// if the user is correct.
//
// Error handling is ignored since email notification service failures are ignored and we still want the user
// to login and perform regular operations even without verification of email.
func PerformEmailNotificationService(user string, emailAddress string, userID string, messageType string) {

	isEmailServiceEnabled := os.Getenv("_SENDGRID_EMAIL_SERVICE")
	if isEmailServiceEnabled != "true" {
		config.Log("email service feature flags are disabled. Email Service is inoperative.", nil)
		return
	}

	WebApplicationEndpoint := os.Getenv("REACT_APP_LOCALHOST_URL")
	if len(WebApplicationEndpoint) <= 0 {
		config.Log("unable to determine the web application endpoint", errors.New(config.ErrorWebApplicationEndpoint))
		return
	}

	// Web App routes are protected and web app uses context to determine password reset. Unable to use JWT to validate user
	// as JWT is only authorized use for logged in users. Web App users pageContext to determine correct navigation as it sits
	// outside of react-router.
	if messageType == config.ResetPasswordTokenStringURI {

		isOtpServiceEnabled := os.Getenv("_OTP_GENERATOR_SERVICE")
		if isOtpServiceEnabled != "true" {
			config.Log("otp service feature flags are disabled. OTP is inoperative.", nil)
			return
		}

		key := os.Getenv("OTP_GENERATOR_API_KEY")
		if len(key) <= 0 {
			config.Log("missing required values", errors.New(config.ErrorGeneratedOTPFailure))
			return
		}

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			config.Log("unable to parse userID", err)
			return
		}

		generatedOTP, err := ether.GenerateOTP(&etherTypes.OTPCredentials{
			UserID:       parsedUserID,
			EmailAddress: emailAddress,
			Token:        key,
		})

		if err != nil {
			config.Log("unable to generate otp token", err)
			return
		}

		// 4. verify token and validity in reset password submission api ( not built )

		err = db.UpdateRecoveryToken(user, userID, generatedOTP)
		if err != nil {
			config.Log("unable to update recovery token", err)
			return
		}

		emailSubjectLine := config.ResetPasswordSubjectLine
		emailPlainTextMessage := config.ResetPasswordPlainTextMessage

		plainText := fmt.Sprintf("Please enter the following digits when prompted: %s", generatedOTP)
		htmlContent := fmt.Sprintf(`<p>%s</p> <h2>%s</h2>`, emailPlainTextMessage, generatedOTP)

		sendMessage(emailAddress, emailSubjectLine, plainText, htmlContent)

	} else {

		generatedToken, err := generateToken(userID)
		if err != nil {
			config.Log("unable to generate token", err)
			return
		}

		messageLinkWithAuthorizationToken := fmt.Sprintf("%s%s?token=%s", WebApplicationEndpoint, config.EmailVerificationTokenStringURI, generatedToken)
		emailSubjectLine := config.VerifyEmailSubjectLine
		emailPlainTextMessage := config.VerifyEmailPlainTextMessage

		plainText := fmt.Sprintf("Please click on the following: %s", generatedToken)
		htmlContent := fmt.Sprintf(`
		<p>%s</p>
		<a href="%s">%s</a>
	`, emailPlainTextMessage, messageLinkWithAuthorizationToken, messageLinkWithAuthorizationToken)

		sendMessage(emailAddress, emailSubjectLine, plainText, htmlContent)
	}

}

// sendMessage ...
//
// function used to send message to the client
func sendMessage(toEmailAddress string, emailSubjectLine string, plainText string, htmlContent string) error {

	sendGridApiKey := os.Getenv("SEND_GRID_API_KEY")
	if len(sendGridApiKey) <= 0 {
		config.Log("email service unavailable", errors.New(config.ErrorInvalidApiKey))
		return errors.New(config.ErrorInvalidApiKey)
	}

	fromUser := os.Getenv("SEND_GRID_USER")
	if len(fromUser) <= 0 {
		config.Log("email service unavailable", errors.New(config.ErrorUserCredentialsNotFound))
		return errors.New(config.ErrorUserCredentialsNotFound)
	}

	fromUserEmailAddress := os.Getenv("SEND_GRID_USER_EMAIL_ADDRESS")
	if len(fromUserEmailAddress) <= 0 {
		config.Log("email service unavailable", errors.New(config.ErrorUserCredentialsNotFound))
		return errors.New(config.ErrorUserCredentialsNotFound)
	}

	from := mail.NewEmail(fromUser, fromUserEmailAddress)
	to := mail.NewEmail(toEmailAddress, toEmailAddress)

	message := mail.NewSingleEmail(from, emailSubjectLine, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(sendGridApiKey)

	_, err := client.Send(message)
	if err != nil {
		config.Log("unable to send email verification", err)
		return errors.New(config.ErrorUnableToSendEmail)
	}

	config.Log("Email notification sent to %s on %+v", nil, toEmailAddress, time.Now())
	return nil
}

// generateToken ...
//
// function is used to generate token to use in the email notification service
func generateToken(userID string) (string, error) {
	formattedTime, err := strconv.ParseInt(config.DefaultTokenValidityTime, 10, 64)
	if err != nil {
		config.Log("unable to parse provided time", err)
		return "", err
	}

	secretToken := os.Getenv("TOKEN_SECRET_KEY")
	if len(secretToken) <= 0 {
		config.Log("unable to retrieve secret token key. defaulting to default values", nil)
		secretToken = ""
	}

	draftCredentials := types.Credentials{
		Claims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: formattedTime,
		},
	}
	credentials, err := stormRider.CreateJWT(&draftCredentials, secretToken)
	if err != nil {
		config.Log("unable to create email token for verification services", err)
		return "", err
	}
	return credentials.Cookie, nil
}

// ValidateCredentials ...
//
// Method is used to verify if the incoming api calls have a valid jwt token.
// If the validity of the token is crossed, or if the token itself is invalid the error is propogated up the method chain.
func ValidateCredentials(user string, ID string) error {
	db, err := db.SetupDB(user)
	if err != nil {
		config.Log("unable to setup db", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		config.Log("unable to setup transaction for selected db", err)
		return err
	}

	var tokenFromDb string
	var expirationTime time.Time
	err = tx.QueryRow(`SELECT token, expiration_time FROM auth.oauth WHERE id=$1 LIMIT 1;`, ID).Scan(&tokenFromDb, &expirationTime)
	if err != nil {
		config.Log("unable to retrive validated token", err)
		tx.Rollback()
		return err
	}

	err = utils.ValidateJwtToken(tokenFromDb)
	if err != nil {
		config.Log("unable to validate jwt token", err)
		tx.Rollback()
		return err
	}

	// Check if the token is within the last 30 seconds of its expiry time
	// token is about to expire. if the user is continuing with activity, create new token
	formattedTimeToLive := time.Until(expirationTime)
	if formattedTimeToLive <= 30*time.Second && formattedTimeToLive > 0 {

		formattedTime, err := strconv.ParseInt(config.DefaultTokenValidityTime, 10, 64)
		if err != nil {
			config.Log("unable to parse provided time", err)
			return err
		}

		secretToken := os.Getenv("TOKEN_SECRET_KEY")
		if len(secretToken) <= 0 {
			config.Log("unable to retrieve secret token key. defaulting to default values", nil)
			secretToken = ""
		}

		draftCredentials := &types.Credentials{
			Claims: jwt.StandardClaims{
				ExpiresAt: formattedTime,
			},
		}

		updatedToken, err := stormRider.RefreshToken(draftCredentials, secretToken)
		if err != nil {
			config.Log("unable to refresh token", err)
			tx.Rollback()
			return err
		}

		parsedUserID, err := uuid.Parse(ID)
		if err != nil {
			config.Log("unable to determine user id", err)
			return err
		}

		tokenValidityMinutes, err := strconv.Atoi(config.DefaultTokenValidityTime)
		if err != nil {
			config.Log("Invalid token validity time: %v", err)
			return err
		}

		draftUser := model.UserCredentials{
			ID:             parsedUserID,
			ExpirationTime: time.Now().Add((time.Duration(tokenValidityMinutes) * time.Minute)),
			PreBuiltToken:  updatedToken,
		}
		err = upsertOauthToken(&draftUser, tx)
		if err != nil {
			config.Log("unable to revalidate the user", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		config.Log("unable to commit transaction", err)
		tx.Rollback()
		return err
	}

	return nil
}
