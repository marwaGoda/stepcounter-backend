package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MarwaAGoda/steps-leaderboard/core"
	"github.com/gorilla/mux"
)

type Handler struct {
	teamManager *core.TeamManager
}

func NewHandler(teamManager *core.TeamManager) *Handler {
	return &Handler{
		teamManager: teamManager,
	}
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Parse the request body and create a new team
	var team core.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}
	team, err = h.teamManager.CreateTeam(team)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, "Failed to create team")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *Handler) GetTeams(w http.ResponseWriter, r *http.Request) {
	// Retrieve a list of all teams from the core package
	// Retrieve a list of all teams from the teamManager
	teams := h.teamManager.GetTeams()
	// Write the response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	// Extract the team ID from the request URL
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	// Retrieve the team details from the core package
	team, err := h.teamManager.GetTeamWithoutIndex(teamID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Team not found")
	}
	// Write the response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
	return
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	// Extract the team ID from the request URL
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	// Parse the request body and create a new user
	var user core.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	users, err := h.teamManager.AddUser(teamID, user)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, "Failed to add user")
		fmt.Println("Error while creating user:", err, user)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Extract the team ID from the request URL
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	users, err := h.teamManager.GetUsers(teamID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, "Failed to get users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	// Extract the team ID and user ID from the request URL
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	userID := vars["userID"]

	err := h.teamManager.IncrementCounter(teamID, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Steps incremented successfully by 1",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetCounter(w http.ResponseWriter, r *http.Request) {
	// Extract the team ID from the request URL
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	counter, err := h.teamManager.GetCounter(teamID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counter)

}

func (h *Handler) IncrementCounterByValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	userID := vars["userID"]

	// Parse the request body and get the increment value
	type IncrementRequest struct {
		IncrementValue int `json:"incrementValue"`
	}
	var req IncrementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	// Call the appropriate methods in the core package
	err = h.teamManager.IncrementCounterByValue(teamID, userID, req.IncrementValue)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]int{
		"counter": req.IncrementValue,
	}
	json.NewEncoder(w).Encode(response)
}
