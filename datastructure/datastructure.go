package datastructure

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post : blog model
type Post struct {
	ID              bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author          string        `json:"author" bson:"author"`
	Title           string        `json:"title" bson:"title"`
	Content         string        `json:"content" bson:"content"`
	IsVerified      bool          `json:"is_verified" bson:"is_verified"`
	LastUpdated     time.Time     `json:"last_updated" bson:"last_updated"`
	CreatedDateTime time.Time     `json:"created_datetime" bson:"created_datetime"`
}

// Comment : comment model
type Comment struct {
	Author          string
	Content         string
	HelpfulCount    int
	CreatedDateTime time.Time
}
