package main

import (
	"log"
	"net/http"

	"github.com/MarwaAGoda/steps-leaderboard/api"
	"github.com/MarwaAGoda/steps-leaderboard/core"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	teamManager *core.TeamManager
}

func NewHandler(teamManager *core.TeamManager) *Handler {
	return &Handler{
		teamManager: teamManager,
	}
}

func main() {
	// Create a new instance of TeamManager
	teamManager := core.NewTeamManager()

	// Create a new instance of Handler and pass the TeamManager
	handler := api.NewHandler(teamManager)

	router := mux.NewRouter()
	// Create a new corsHandler to enable CORS for all routes
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),            // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST"}),  // Allow GET and POST methods
		handlers.AllowedHeaders([]string{"Content-Type"}), // Allow only "Content-Type" header
	)

	// Register API routes using the Handler methods
	router.HandleFunc("/teams", handler.CreateTeam).Methods("POST")
	router.HandleFunc("/teams", handler.GetTeams).Methods("GET")
	router.HandleFunc("/teams/{teamID}", handler.GetTeam).Methods("GET")
	router.HandleFunc("/teams/{teamID}/users", handler.AddUser).Methods("POST")
	router.HandleFunc("/teams/{teamID}/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/teams/{teamID}/users/{userID}/counters", handler.IncrementCounter).Methods("POST")
	router.HandleFunc("/teams/{teamID}/counters", handler.GetCounter).Methods("GET")
	router.HandleFunc("/teams/{teamID}/users/{userID}/counters/increment", handler.IncrementCounterByValue).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}
