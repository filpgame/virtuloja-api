package models

import "gopkg.in/mgo.v2/bson"

type (
	// Alert represents the structure of our resource
	Alert struct {
		ID    bson.ObjectId
		Alert string
	}
)
