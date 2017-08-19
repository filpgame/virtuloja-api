package models

import "gopkg.in/mgo.v2/bson"

type (
	// Product represents the structure of our resource
	Product struct {
		ID          bson.ObjectId
		GlobalID    int32
		Value       float32
		Description string
	}
)
