package core

import (
	"errors"
	"sync"
)

type TeamManager struct {
	teams map[string]Team
	mu    sync.Mutex
}

func NewTeamManager() *TeamManager {
	return &TeamManager{
		teams: make(map[string]Team),
	}
}

func (tm *TeamManager) CreateTeam(team Team) (Team, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// Check if team name is unique
	for _, t := range tm.teams {
		if t.Name == team.Name {
			return Team{}, errors.New("team name already exists")
		}
	}
	team = NewTeam(team.Name)
	if _, exists := tm.teams[team.ID]; exists {
		return Team{}, errors.New("team with the same ID already exists")
	}

	tm.teams[team.ID] = team
	return team, nil
}

func (tm *TeamManager) GetTeams() []Team {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	teams := make([]Team, 0, len(tm.teams))
	for _, team := range tm.teams {
		teams = append(teams, team)
	}
	return teams
}

func (tm *TeamManager) GetTeam(teamID string) (Team, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return Team{}, errors.New("team not found")
	}
	return team, nil
}

func (tm *TeamManager) GetTeamWithoutIndex(teamID string) (Team, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return Team{}, errors.New("team not found")
	}
	return team, nil
}

func (tm *TeamManager) AddUser(teamID string, user User) ([]User, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return nil, errors.New("team not found")
	}
	// Check if user name is unique within the team
	for _, u := range team.Users {
		if u.Name == user.Name {
			return nil, errors.New("user name already exists in the team")
		}
	}
	// Add user to the team's users list
	team.Users = append(team.Users, NewUser(user.Name, user.Counter))
	tm.teams[teamID] = team

	return team.Users, nil
}

func (tm *TeamManager) GetUsers(teamID string) ([]User, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return nil, errors.New("team not found")
	}

	return team.Users, nil
}

func (tm *TeamManager) GetUserWithIndex(teamID, userID string) (User, int, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return User{}, -1, errors.New("team not found")
	}

	for i, user := range team.Users {
		if user.ID == userID {
			return user, i, nil
		}
	}

	return User{}, -1, errors.New("user not found")
}

func (tm *TeamManager) GetUserIndex(teamID string, userID string) (int, error) {
	team, exists := tm.teams[teamID]
	if !exists {
		return -1, errors.New("team not found")
	}
	// get user of the team
	for index, user := range team.Users {
		if user.ID == userID {
			return index, nil
		}
	}

	return -1, errors.New("user not found")

}

func (tm *TeamManager) IncrementCounter(teamID, userID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return errors.New("team not found")
	}

	for i, user := range team.Users {
		if user.ID == userID {
			team.Users[i].Counter++
			tm.teams[teamID] = team
			team.Counter.Count++
			return nil
		}
	}

	return errors.New("user not found")
}

func (tm *TeamManager) GetCounter(teamID string) (int, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return 0, errors.New("team not found")
	}

	counter := 0
	for _, user := range team.Users {
		counter += user.Counter
	}

	return counter, nil
}

func (tm *TeamManager) IncrementCounterByValue(teamID, userID string, value int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return errors.New("team not found")
	}
	index, error := tm.GetUserIndex(teamID, userID)
	if error != nil {
		return errors.New("user not found")
	}
	team.Counter.Count += value
	team.Users[index].Counter += value
	tm.teams[team.ID] = team
	return nil
}
