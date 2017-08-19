package models

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	// Cart represents the structure of our resource
	Cart struct {
		ID       bson.ObjectId
		Products int32
	}
)
