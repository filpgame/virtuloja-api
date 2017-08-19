package main

import (
	// Standard library packages

	"net/http"

	// Third party packages
	"github.com/filpgame/virtuloja-api/controllers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {

	// Instantiate a new router
	r := httprouter.New()

	// Get a UserController instance
	pc := controllers.NewProductController(getSession())

	// Get a product resource
	r.GET("/product/:id", pc.GetProduct)

	// Get a product resource
	r.GET("/productGlobalId/:id", pc.GetProductByGlobalID)

	// Create a new product
	r.POST("/product", pc.CreateProduct)

	// Get a UserController instance
	sc := controllers.NewStockController(getSession())

	// Create a new product
	r.POST("/stock", sc.RegisterStockItem)

	// Remove an existing product
	// r.DELETE("/product/:id", uc.RemoveUser)

	// Get a UserController instance
	cc := controllers.NewCartController(getSession())

	// Create a new product
	r.POST("/checkin", cc.CreateCart)

	// Create a new product
	r.POST("/checkout/:customerId", cc.Checkout)

	// Create a new product
	r.POST("/cart/:customerId/addProduct", cc.AddProduct)

	// Fire up the server
	http.ListenAndServe(":8000", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
