// models.user.go

package springkilometers

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)

type user struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"-"`
	Salt     string `json:"-"`
}

// random salt with given length
func salting(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// sha3 hash algorithm
func crypting(s string) string {
	h := sha3.New512()
	h.Write([]byte(s))
	pass := h.Sum(nil)
	return string(pass)
}

// Check if the username and password combination is valid
func isUserValid(username, password string) bool {

	if isUsernameAvailable(username) {
		return false
	}

	// TODO: Get salt, hashed password from existing user
	// SQL: SELECT users.password, users.salt FROM users WHERE users.username = $username -> salt, cryptpass
	cryptpass := "tmp"
	dbsalt := "tmp"

	// TODO: hash function + salt to input password
	pass := crypting(password + dbsalt)

	// TODO: compare username, password(hashed) with input

	return pass == cryptpass
}

// Register a new user with the given username and password
func RegisterNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	salt := salting(10)
	pass := crypting(password + salt)

	u := user{Username: username, Password: pass, Salt: salt}

	// TODO: add $user to database
	// SQL: INSERT INTO users ( username, password, salt) VALUES ( $username, $pass, $salt)

	return &u, nil
	// idk if return
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	// TODO: exist sql command
	// SQL: SELECT users.username FROM users where users.username = $username
	return true
}
