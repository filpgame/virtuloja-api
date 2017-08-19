package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/swhite24/go-rest-tutorial/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Product represents the structure of our resource
	Product struct {
		ID          bson.ObjectId
		GlobalID    int32
		Value       float32
		Description string
	}
)

type (
	// ProductController represents the controller for operating on the User resource
	ProductController struct {
		session *mgo.Session
	}
)

// NewProductController provides a reference to a UserController with provided mongo session
func NewProductController(s *mgo.Session) *ProductController {
	return &ProductController{s}
}

// GetProduct retrieves an individual user resource
func (pc ProductController) GetProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub product
	u := models.User{}

	// Fetch product
	if err := pc.session.DB("virtuloja-api").C("products").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateProduct creates a new user resource
func (pc ProductController) CreateProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Stub a product to be populated from the body
	pr := Product{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&pr)

	// Add an Id
	pr.ID = bson.NewObjectId()

	// Write the product to mongo
	pc.session.DB("virtuloja-api").C("products").Insert(pr)

	// Marshal provided interface into JSON structure
	prj, _ := json.Marshal(pr)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", prj)
}

// RemoveProduct removes an existing product resource
func (pc ProductController) RemoveProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove product
	if err := pc.session.DB("virtuloja-api").C("products").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
