package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/filpgame/virtuloja-api/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// CartController represents the controller for operating on the Cart resource
	CartController struct {
		session *mgo.Session
	}
)

// NewCartController provides a reference to a UserController with provided mongo session
func NewCartController(s *mgo.Session) *CartController {
	return &CartController{s}
}

// // GetCart retrieves an individual cartresource by ID
// func (pc CartController) GetCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	// Grab id
// 	id := p.ByName("id")

// 	// Verify id is ObjectId, otherwise bail
// 	if !bson.IsObjectIdHex(id) {
// 		msg := "Invalid ObjectId"
// 		log.Println(msg)
// 		http.Error(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	// Grab id
// 	oid := bson.ObjectIdHex(id)

// 	// Stub product
// 	pr := models.Product{}

// 	// Fetch product
// 	if err := pc.session.DB("virtuloja-api").C("products").FindId(oid).One(&pr); err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Marshal provided interface into JSON structure
// 	uj, _ := json.Marshal(pr)

// 	// Write content-type, statuscode, payload
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	fmt.Fprintf(w, "%s", uj)
// }

// GetCartByCustomerID retrieves an individual product resource by global id
// func (pc CartController) GetCartByCustomerID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	// Grab id
// 	id, err := strconv.Atoi(p.ByName("id"))
// 	if err != nil {
// 		msg := "Invalid id"
// 		log.Println(msg)
// 		http.Error(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	// Query All
// 	var results []models.Product
// 	err = pc.session.DB("virtuloja-api").C("products").Find(bson.M{"globalid": id}).Sort("-timestamp").All(&results)

// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if len(results) == 0 {
// 		msg := "No results found"
// 		log.Println(msg)
// 		http.Error(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	// Marshal provided interface into JSON structure
// 	uj, _ := json.Marshal(results)

// 	// Write content-type, statuscode, payload
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	fmt.Fprintf(w, "%s", uj)
// }

// CreateCart creates a new user resource
func (cc CartController) CreateCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Stub a product to be populated from the body
	cAux := models.Cart{}
	c := models.Cart{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&cAux)
	c.CustomerID = cAux.CustomerID
	c.TimeIni = time.Now()

	// Add an Id
	c.ID = bson.NewObjectId()

	// Write the product to mongo
	cc.session.DB("virtuloja-api").C("carts").Insert(c)

	// Marshal provided interface into JSON structure
	prj, _ := json.Marshal(c)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", prj)
}

// AddProduct creates a new user resource
func (cc CartController) AddProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Grab id
	customerID, err := strconv.Atoi(p.ByName("customerId"))

	// Stub a product to be populated from the body
	i := models.Item{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&i)

	var c models.Cart

	err = cc.session.DB("virtuloja-api").C("carts").Find(bson.M{"customerid": customerID}).One(&c)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Adicionar novos itens na lista
	repeated := false
	for k := 0; k < len(c.Items); k++ {

		if i.Product == c.Items[k].Product {

			c.Items[k].Quantity += i.Quantity
			repeated = true
			break
		}
	}

	if !repeated {

		c.Items = append(c.Items, i)
	}

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"items": c.Items}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	info, err := cc.session.DB("virtuloja-api").C("carts").Find(bson.M{"customerid": c.CustomerID}).Apply(change, &c)
	fmt.Println(info)
	fmt.Println(c.CustomerID)

	// Add an Id
	c.ID = bson.NewObjectId()

	// Write the product to mongo
	//cc.session.DB("virtuloja-api").C("carts").Insert(c)

	// Marshal provided interface into JSON structure
	prj, _ := json.Marshal(c)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", prj)
}

// RemoveProduct removes an existing product resource
// func (pc CartController) RemoveProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

// 	// Grab id
// 	id := p.ByName("id")

// 	// Verify id is ObjectId, otherwise bail
// 	if !bson.IsObjectIdHex(id) {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	// Grab id
// 	oid := bson.ObjectIdHex(id)

// 	// Remove product
// 	if err := pc.session.DB("virtuloja-api").C("products").RemoveId(oid); err != nil {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	// Write status
// 	w.WriteHeader(200)
// }
