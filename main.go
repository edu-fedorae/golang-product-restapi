package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/api/v1/products", createProduct).Methods("POST")
	r.HandleFunc("/api/v1/products", getProducts).Methods("GET")
	r.HandleFunc("/api/v1/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/v1/products/{id}", updateProduct).Methods("PATCH")
	r.HandleFunc("/api/v1/products/{id}", deleteProduct).Methods("DELETE")

	// db := connect()
	// defer db.Close()

	// product := Product{ID: uuid.New(), Name: "My Product", Quantity: 10, Price: 49.99}
	// fmt.Println(product)

	// Starting Server
	log.Fatal(http.ListenAndServe(":8000", r))
}

// homePage
func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

// Create Product
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Creating Product Instance
	product := &Product{
		ID: uuid.New().String(),
	}

	// Decoding Request
	_ = json.NewDecoder(r.Body).Decode(&product)

	// Inserting Into Database
	_, err := db.Model(product).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Product
	json.NewEncoder(w).Encode(product)
}

// Get Products
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Creating Products Slice
	var products []Product
	if err := db.Model(&products).Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Products
	json.NewEncoder(w).Encode(products)
}

// Get Product
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	// Creating Product Instance
	product := &Product{ID: productId}
	if err := db.Model(product).WherePK().Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Product
	json.NewEncoder(w).Encode(product)
}

// Update Product & Store*
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	// Creating Product Instance
	product := &Product{ID: productId}

	_ = json.NewDecoder(r.Body).Decode(&product)

	// Alternate Way To Include Store In Update
	// store := map[string]interface{}{
	// 	"id": product.Store.ID,
	// 	"name": product.Store.Name,
	// }
	// db.Model(product).WherePK().Set("name = ?, quantity = ?, price = ?, store = ?", product.Name, product.Quantity, product.Price, store).Update()

	_, err := db.Model(product).WherePK().Set("name = ?, quantity = ?, price = ?, store = ?", product.Name, product.Quantity, product.Price, product.Store).Update()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Product
	json.NewEncoder(w).Encode(product)
}

// Delete Product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	// Creating Product Instance Alternative Way
	// product := &Product{ID: productId}
	// result, err := db.Model(product).WherePK().Delete()

	// Creating Product Instance
	product := &Product{}
	result, err := db.Model(product).Where("id = ?", productId).Delete()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning result
	json.NewEncoder(w).Encode(result)
}
