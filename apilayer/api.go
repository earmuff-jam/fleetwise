// # Documentation for the Fleet wise api layer.
//
// Scehmes: https
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/earmuff-jam/fleetwise/bucket"
	"github.com/earmuff-jam/fleetwise/config"
	"github.com/earmuff-jam/fleetwise/handler"
	"github.com/earmuff-jam/fleetwise/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// MessageResponse ...
// swagger:model MessageResponse
//
// MessageResponse struct
type MessageResponse struct {
	Message string
}

// CustomRequestHandler function
//
// wrapper function for all request, response pair
type CustomRequestHandler func(http.ResponseWriter, *http.Request, string)

func main() {

	// load environment variables
	instance := os.Getenv("ENVIRONMENT")
	if instance == "" {
		err := godotenv.Load(filepath.Join("..", ".env"))
		if err != nil {
			log.Printf("No env file detected. Using os system configuration.")
		}
	}

	config.InitLogger()
	bucket.InitializeStorageAndBucket()

	router := mux.NewRouter()

	// public routes
	router.HandleFunc("/api/v1/signup", handler.Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/verify", handler.VerifyEmailAddress).Methods("GET")

	router.HandleFunc("/api/v1/signin", handler.Signin).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/isValidEmail", handler.IsValidUserEmail).Methods("POST")
	router.HandleFunc("/api/v1/resetPassword", handler.ResetPassword).Methods("POST")
	router.HandleFunc("/api/v1/logout", handler.Logout).Methods("GET", "OPTIONS")

	// secure routes
	router.Handle("/api/v1/reset", CustomRequestHandler(handler.ResetEmailToken)).Methods(http.MethodPost)
	router.Handle("/api/v1/locations", CustomRequestHandler(handler.GetAllStorageLocations)).Methods(http.MethodGet)

	// summary
	router.Handle("/api/v1/summary", CustomRequestHandler(handler.GetAssetsAndSummary)).Methods(http.MethodGet)

	// categories
	router.Handle("/api/v1/category/items", CustomRequestHandler(handler.GetAllCategoryItems)).Methods(http.MethodGet)
	router.Handle("/api/v1/category/items", CustomRequestHandler(handler.AddItemsInCategory)).Methods(http.MethodPost)
	router.Handle("/api/v1/category/remove/items", CustomRequestHandler(handler.RemoveAssociationFromCategory)).Methods(http.MethodPost)

	router.Handle("/api/v1/categories", CustomRequestHandler(handler.GetAllCategories)).Methods(http.MethodGet)
	router.Handle("/api/v1/category", CustomRequestHandler(handler.GetCategory)).Methods(http.MethodGet)
	router.Handle("/api/v1/category", CustomRequestHandler(handler.CreateCategory)).Methods(http.MethodPost)
	router.Handle("/api/v1/category/{id}", CustomRequestHandler(handler.UpdateCategory)).Methods(http.MethodPut)
	router.Handle("/api/v1/category/{id}", CustomRequestHandler(handler.RemoveCategory)).Methods(http.MethodDelete)

	// maintenance plans
	router.Handle("/api/v1/plans/items", CustomRequestHandler(handler.GetAllMaintenancePlanItems)).Methods(http.MethodGet)
	router.Handle("/api/v1/plans/items", CustomRequestHandler(handler.AddItemsInMaintenancePlan)).Methods(http.MethodPost)
	router.Handle("/api/v1/plan/remove/items", CustomRequestHandler(handler.RemoveAssociationFromMaintenancePlan)).Methods(http.MethodPost)

	router.Handle("/api/v1/maintenance-plans", CustomRequestHandler(handler.GetAllMaintenancePlans)).Methods(http.MethodGet)
	router.Handle("/api/v1/plan", CustomRequestHandler(handler.CreateMaintenancePlan)).Methods(http.MethodPost)
	router.Handle("/api/v1/plan", CustomRequestHandler(handler.GetMaintenancePlan)).Methods(http.MethodGet)
	router.Handle("/api/v1/plan/{id}", CustomRequestHandler(handler.UpdateMaintenancePlan)).Methods(http.MethodPut)
	router.Handle("/api/v1/plan/{id}", CustomRequestHandler(handler.RemoveMaintenancePlan)).Methods(http.MethodDelete)

	// profile
	router.Handle("/api/v1/profile/list", CustomRequestHandler(handler.GetAllUserProfiles)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}", CustomRequestHandler(handler.GetProfile)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/stats", CustomRequestHandler(handler.GetProfileStats)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/notifications", CustomRequestHandler(handler.GetNotifications)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/notifications", CustomRequestHandler(handler.UpdateSelectedMaintenanceNotification)).Methods(http.MethodPut)
	router.Handle("/api/v1/profile/{id}/recent-activities", CustomRequestHandler(handler.GetRecentActivities)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}", CustomRequestHandler(handler.UpdateProfile)).Methods(http.MethodPut)
	router.Handle("/api/v1/profile/{id}/username", CustomRequestHandler(handler.GetUsername)).Methods(http.MethodGet)

	// inventories
	router.Handle("/api/v1/profile/{id}/inventories", CustomRequestHandler(handler.GetAllInventories)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/inventories/{invID}", CustomRequestHandler(handler.GetInventoryByID)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/inventories/{asssetID}", CustomRequestHandler(handler.UpdateAssetColumn)).Methods(http.MethodPut)

	router.Handle("/api/v1/profile/{id}/inventories", CustomRequestHandler(handler.AddNewInventory)).Methods(http.MethodPost)
	router.Handle("/api/v1/profile/{id}/inventories/bulk", CustomRequestHandler(handler.AddInventoryInBulk)).Methods(http.MethodPost)
	router.Handle("/api/v1/profile/{id}/inventories", CustomRequestHandler(handler.UpdateSelectedInventory)).Methods(http.MethodPut)
	router.Handle("/api/v1/profile/{id}/inventories/prune", CustomRequestHandler(handler.RemoveSelectedInventory)).Methods(http.MethodPost)

	router.Handle("/api/v1/profile/{id}/fav", CustomRequestHandler(handler.GetFavouriteItems)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/fav", CustomRequestHandler(handler.SaveFavItem)).Methods(http.MethodPost)
	router.Handle("/api/v1/profile/{id}/fav", CustomRequestHandler(handler.RemoveFavItem)).Methods(http.MethodDelete)

	// notes
	router.Handle("/api/v1/profile/{id}/notes", CustomRequestHandler(handler.GetNotes)).Methods(http.MethodGet)
	router.Handle("/api/v1/profile/{id}/notes", CustomRequestHandler(handler.AddNewNote)).Methods(http.MethodPost)
	router.Handle("/api/v1/profile/{id}/notes", CustomRequestHandler(handler.UpdateNote)).Methods(http.MethodPut)
	router.Handle("/api/v1/profile/{id}/notes/{noteID}", CustomRequestHandler(handler.RemoveNote)).Methods(http.MethodDelete)

	// reports
	router.Handle("/api/v1/reports/{id}", CustomRequestHandler(handler.GetReports)).Methods(http.MethodGet)

	// image
	router.Handle("/api/v1/{id}/uploadImage", CustomRequestHandler(handler.UploadImage)).Methods(http.MethodPost)
	router.Handle("/api/v1/{id}/fetchImage", CustomRequestHandler(handler.FetchImage)).Methods(http.MethodGet)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Role2"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}),
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"http://localhost", "http://localhost:5173", "http://localhost:5173", "http://localhost:8081"}),
		handlers.ExposedHeaders([]string{"Role2"}),
	)

	http.Handle("/", cors(router))

	config.Log("Api is up and running ...", nil)
	err := http.ListenAndServe(":8087", nil)
	if err != nil {
		config.Log("failed to start the server", err)
		return
	}
}

// ServerHTTP is a wrapper function to derieve the user authentication.
//
// Performs check to validate if there is a current user in the system that
// can communication with the db as a user and also validates jwt for incoming
// requests and refresh token if necessary.
func (u CustomRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		config.Log("missing license key", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentUser := validateCurrentUser()
	err = service.ValidateCredentials(currentUser, cookie.Value)
	if err != nil {
		config.Log("failed to validate token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u(w, r, currentUser)
}

// validateCurrentUser ...
//
// function is used to check if currentUser exists in the system. defaults to system user from os/user
func validateCurrentUser() string {
	currentUser := os.Getenv("CLIENT_USER")
	if len(currentUser) == 0 {
		user, _ := user.Current()
		config.Log("unable to retrieve user from env. Using user - %s", nil, user.Username)
		currentUser = user.Username
	}
	return currentUser
}
