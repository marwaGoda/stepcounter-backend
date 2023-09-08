package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MarwaAGoda/steps-leaderboard/core"
	"github.com/gorilla/mux"
)

func TestCreateTeam(t *testing.T) {
	tm := core.NewTeamManager()
	handler := NewHandler(tm)
	router := mux.NewRouter()
	router.HandleFunc("/teams", handler.CreateTeam).Methods("POST")

	// Test Case: Creating a new team successfully
	payload := `{"name":"Team A"}`
	req, _ := http.NewRequest("POST", "/teams", bytes.NewBufferString(payload))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("Expected status code %v but got %v", http.StatusCreated, res.Code)
	}

}

func TestAddUser(t *testing.T) {
	tm := core.NewTeamManager()
	handler := NewHandler(tm)

	// Create a new team
	newTeam := core.Team{
		Name: "Team A",
	}
	newTeam, _ = tm.CreateTeam(newTeam)

	// Create a request with a JSON payload for a new user
	payload := `{"name":"User A"}`
	req, err := http.NewRequest("POST", "/teams/"+newTeam.ID+"/users", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/teams/{teamID}/users", handler.AddUser).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}
	fmt.Printf("t: %v\n", rr.Body)
	// Check the response body
	var users []core.User
	err = json.NewDecoder(rr.Body).Decode(&users)
	if err != nil {
		t.Fatalf("Failed to parse response: %v\nResponse body: %s", err, rr.Body.String())
	}

	// Validate the response
	if len(users) != 1 {
		t.Fatalf("Expected 1 user in the response, but got %d", len(users))
	}

	// Check the user data
	expectedUser := core.User{ID: users[0].ID, Name: "User A"}
	if !reflect.DeepEqual(users[0], expectedUser) {
		t.Errorf("Expected user %+v, but got %+v", expectedUser, users[0])
	}
}

// Test for IncrementCounter handler
func TestIncrementCounter(t *testing.T) {
	tm := core.NewTeamManager()
	handler := NewHandler(tm)

	// Create a new team and add a user to it
	newTeam := core.Team{
		Name: "Team A",
	}
	newTeam, _ = tm.CreateTeam(newTeam)
	var users []core.User
	users, _ = tm.AddUser(newTeam.ID, core.User{Name: "User A"})

	// Create a request to increment the counter of a specific user in the team
	req, err := http.NewRequest("POST", "/teams/"+newTeam.ID+"/users/"+users[0].ID+"/counters", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/teams/{teamID}/users/{userID}/counters", handler.IncrementCounter).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Ensure the step counter has been incremented
	team, err := tm.GetTeam(newTeam.ID)
	user, _, err := tm.GetUserWithIndex(team.ID, users[0].ID)

	expectedCounter := 1
	if user.Counter != expectedCounter {
		t.Errorf("Expected step counter to be %d, but got %d", expectedCounter, user.Counter)
	}
}

func TestGetCounter(t *testing.T) {
	tm := core.NewTeamManager()
	handler := NewHandler(tm)

	newTeam := core.Team{
		Name: "Team A",
	}
	newTeam, _ = tm.CreateTeam(newTeam)

	// Create a request to get the current step counter of a team
	req, err := http.NewRequest("GET", "/teams/"+newTeam.ID+"/counters", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/teams/{teamID}/counters", handler.GetCounter).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	var counter int
	err = json.NewDecoder(rr.Body).Decode(&counter)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

}

func TestIncrementCounterByValue(t *testing.T) {
	tm := core.NewTeamManager()
	handler := NewHandler(tm)

	// Create a new team and add a user to it
	newTeam := core.Team{
		Name: "Team A",
	}
	newTeam, _ = tm.CreateTeam(newTeam)
	var users []core.User
	users, _ = tm.AddUser(newTeam.ID, core.User{Name: "User A"})

	// Create a request to increment the counter of a specific user in the team by a value
	incrementValue := 5
	payload := fmt.Sprintf(`{"incrementValue": %d}`, incrementValue)
	req, err := http.NewRequest("POST", "/teams/"+newTeam.ID+"/users/"+users[0].ID+"/counters/increment", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/teams/{teamID}/users/{userID}/counters/increment", handler.IncrementCounterByValue).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	team, _ := tm.GetTeam(newTeam.ID)
	user, _, err := tm.GetUserWithIndex(team.ID, users[0].ID)
	if err != nil {
		t.Fatal("User not found in team")
	}

	expectedCounter := 5
	if *&user.Counter != expectedCounter {
		t.Errorf("Expected step counter to be %d, but got %d", expectedCounter, *&user.Counter)
	}
	if *&team.Counter.Count != expectedCounter {
		t.Errorf("Expected step counter to be %d, but got %d", expectedCounter, *&team.Counter.Count)
	}

}
