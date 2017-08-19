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

// CreateCart creates a new user resource
func (cc CartController) CreateCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Stub a product to be populated from the body
	cAux := models.Cart{}
	c := models.Cart{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&cAux)
	c.CustomerID = cAux.CustomerID
	c.TimeStart = time.Now()

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
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Stub a product to be populated from the body
	i := models.Item{}

	// Populate the product data
	json.NewDecoder(r.Body).Decode(&i)

	var pp models.Product
	err = cc.session.DB("virtuloja-api").C("products").Find(bson.M{"globalid": i.GlobalID}).One(&pp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i.Value = pp.Value

	// Stock validation
	var st models.StockItem
	if err := cc.session.DB("virtuloja-api").C("stock").Find(bson.M{"globalid": i.GlobalID}).One(&st); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if st.Quantity < i.Quantity {

		msg := "There are no enough products " + pp.Description
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

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

		if i.GlobalID == c.Items[k].GlobalID {

			c.Items[k].Quantity += i.Quantity
			repeated = true
			break
		}
	}

	if !repeated {

		c.Items = append(c.Items, i)
	}

	// Update customer cart
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"items": c.Items}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}

	_, err = cc.session.DB("virtuloja-api").C("carts").Find(bson.M{"customerid": c.CustomerID}).Apply(change, &c)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stockQuantityBefore := st.Quantity

	// Update stock
	change = mgo.Change{
		Update:    bson.M{"$set": bson.M{"quantity": st.Quantity - i.Quantity}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}

	_, err = cc.session.DB("virtuloja-api").C("stock").Find(bson.M{"globalid": i.GlobalID}).Apply(change, &st)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if st.Quantity < st.MinimumStock && stockQuantityBefore >= st.MinimumStock {

		// Create low stock alert
		var a models.Alert
		a.ID = bson.NewObjectId()
		a.Alert = "Low stock from product " + pp.Description
		err = cc.session.DB("virtuloja-api").C("alert").Insert(a)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	c.Sum = models.Sum(c)

	// Marshal provided interface into JSON structure
	prj, _ := json.Marshal(c)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", prj)
}

// Checkout retrieves an individual cart resource by ID
func (cc CartController) Checkout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Grab id
	customerID, err := strconv.Atoi(p.ByName("customerId"))

	// Stub product
	c := models.Cart{}

	// Fetch product
	if err := cc.session.DB("virtuloja-api").C("carts").Find(bson.M{"customerid": customerID}).One(&c); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"timeend": time.Now()}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	_, err = cc.session.DB("virtuloja-api").C("carts").Find(bson.M{"customerid": c.CustomerID}).Apply(change, &c)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Sum = models.Sum(c)

	// Write the chart to history
	cc.session.DB("virtuloja-api").C("customerHistory").Insert(c)

	//
	err = cc.session.DB("virtuloja-api").C("carts").Remove(bson.M{"customerid": c.CustomerID})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(c)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
