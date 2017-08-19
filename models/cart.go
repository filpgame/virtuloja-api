package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Item represents the structure of our resource
	Item struct {
		Value        float32
		GlobalID     int32
		Quantity     int16
		MinimumStock int16
	}
)

type (
	// Cart represents the structure of our resource
	Cart struct {
		ID         bson.ObjectId
		CustomerID int32
		Items      []Item
		TimeStart  time.Time
		TimeEnd    time.Time
		Sum        float32
	}
)

// Sum calculates the sum of the values of the cart
func Sum(c Cart) float32 {

	var sum float32
	sum = 0

	for j := 0; j < len(c.Items); j++ {

		sum += float32(c.Items[j].Quantity) * c.Items[j].Value
	}

	return sum
}
