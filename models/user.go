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

// User --
type User struct {
	gorm.Model
	UserID   int    `"json:"userid`
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

// IsUserValid --
// Check if the username and password combination is valid
func IsUserValid(username, password string) bool {

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

// RegisterNewUser a new user with the given username and password
func RegisterNewUser(username, password string) (*User, error) {
	var user User
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	salt := salting(10)
	pass := crypting(password + salt)

	user = User{Username: username, Password: pass, Salt: salt}
	result := db.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	var user User
	db.Where("name = ?", username).First(&user)
	if user.UserID > 0 {
		return false
	}

	return true
}
