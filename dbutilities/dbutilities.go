package dbutilities

import (
	"log"

	"github.com/tanmaybaranwal/BlogAppUsingMGO/datastructure"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Connect : Establish a connection to MongoDB
func Connect(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

// Insert : inserts the object into the collection
func Insert(c *mgo.Collection, p datastructure.Post) (string, bool) {
	var oID = bson.NewObjectId()
	p.ID = oID
	err := c.Insert(&p)
	if err != nil {
		log.Println("Error in the Insert function: ", err)
		return "", false
	}
	return string(p.ID.Hex()), true
}

// FindAll : Get all the data from the collection
func FindAll(c *mgo.Collection, sortBy []string) []datastructure.Post {
	var results []datastructure.Post
	err := c.Find(nil).Sort(sortBy...).All(&results)
	if err != nil {
		log.Println("Error in the FindAll function: ", err)
	}
	return results
}

// FindQuery : Get all documents based on query
func FindQuery(c *mgo.Collection, query bson.M, sortBy []string) []datastructure.Post {
	var results []datastructure.Post
	err := c.Find(query).Sort(sortBy...).All(&results)
	if err != nil {
		log.Println("Error in the FindQuery function: ", err)
	}
	return results
}

// Find : Search for an element
func Find(c *mgo.Collection, query bson.M) datastructure.Post {
	var result datastructure.Post
	err := c.Find(query).One(&result)
	if err != nil {
		log.Println("Error in the Find function: ", err)
	}
	return result
}
