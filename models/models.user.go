// models.user.go

package main

import (
	"errors"
	"strings"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

// Moreover, in production settings, you should
// store passwords securely by salting and hashing them instead
// of using them as we're doing in this demo
var userList = []user{
	user{Username: "user1", Password: "pass1"},
	user{Username: "user2", Password: "pass2"},
	user{Username: "user3", Password: "pass3"},
}

// Check if the username and password combination is valid
func isUserValid(username, password string) bool {
	// TODO: hash function + salt

	// TODO: compare username, password(hashed) with input
}

// Register a new user with the given username and password
func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}
	userList = append(userList, u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	// TODO: exist sql command
	// SQL: SELECT users.username FROM users where users.username = $username
}
