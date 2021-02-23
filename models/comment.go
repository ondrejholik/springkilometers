package springkilometers

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/enescakir/emoji"
)

// Comment model
type Comment struct {
	UserID    int    `json:"user_id"`
	TripID    int    `json:"trip_id"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

// CommentResult --
type CommentResult struct {
	UserID    int    `json:"user_id"`
	TripID    int    `json:"trip_id"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Avatar    string `json:"avatar"`
	Username  string `json:"username"`
}

// AddComment --
func AddComment(message []byte, userID int, room string) {

	go MyCache.Delete(Ctx, "trip:"+room)
	var comment Comment
	var err error
	comment.Timestamp = int(time.Now().Unix())
	json.Unmarshal(message, &comment)
	comment.Message = emoji.Parse(comment.Message)

	comment.UserID = userID
	comment.TripID, err = strconv.Atoi(room)

	if err == nil {
		db.Create(&comment)
	}
}

// GetComments --
func GetComments(tripID int) []CommentResult {
	var comments []CommentResult
	db.Raw("SELECT * FROM comments LEFT JOIN users on users.id = comments.user_id WHERE trip_id = ? ORDER BY timestamp", tripID).Scan(&comments)
	log.Printf("%+v\n", comments)
	return comments
}
