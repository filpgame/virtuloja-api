package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Item represents the structure of our resource
	Item struct {
		Product  int32
		Quantity int16
	}
)

type (
	// Cart represents the structure of our resource
	Cart struct {
		ID         bson.ObjectId
		CustomerID int32
		Items      []Item
		TimeIni    time.Time
		TimeFim    time.Time
	}
)
