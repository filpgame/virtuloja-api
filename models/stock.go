package models

import "gopkg.in/mgo.v2/bson"

type (
	// StockItem represents the structure of our resource
	StockItem struct {
		ID           bson.ObjectId
		GlobalID     int32
		Value        float32
		Description  string
		Quantity     int16
		MinimumStock int16
	}
)
