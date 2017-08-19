package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/filpgame/virtuloja-api/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// StockController represents the controller for operating on the User resource
	StockController struct {
		session *mgo.Session
	}
)

// NewStockController provides a reference to a UserController with provided mongo session
func NewStockController(s *mgo.Session) *StockController {
	return &StockController{s}
}

// RegisterStockItem creates a new user resource
func (sc StockController) RegisterStockItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Stub a product to be populated from the body
	s := models.StockItem{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&s)

	// Add an Id
	s.ID = bson.NewObjectId()

	// Write the product to mongo
	sc.session.DB("virtuloja-api").C("stock").Insert(s)

	// Marshal provided interface into JSON structure
	prj, _ := json.Marshal(s)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", prj)
}
