package core

import (
	"github.com/google/uuid"
)

type Team struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Users   []User  `json:"users"`
	Counter Counter `json:"counter"`
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Counter int    `json:"counter"`
}

type Counter struct {
	Count int `json:"count"`
}

// NewTeam creates a new team with a generated ID.
func NewTeam(name string) Team {
	return Team{
		ID:      uuid.New().String(),
		Name:    name,
		Users:   make([]User, 0),
		Counter: Counter{},
	}
}

func NewUser(name string, counter int) User {
	return User{
		ID:      uuid.New().String(), // Generate a unique ID using UUID
		Name:    name,
		Counter: counter,
	}
}
