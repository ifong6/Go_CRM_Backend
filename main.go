package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Customer struct {
	ID uint8
	Name string
	Role string
	Email string
	Phone string
	Contacted bool//(i.e., indication of whether or not the customer has been contacted)
}

// map - mock "database" to store customer data
var	db = map[uint8]Customer{}

// -----------------------------------------------------------
// ---------------------------APIS----------------------------
func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get the "id" from the URL path
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Use the helper function to get the customer
	customer, success := isValidId(w, idStr)
	if !success {
		return // The error is already handled in the helper function
	}

	// If the customer exists, encode the customer details as JSON and send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a new customer instance
    var newCustomer Customer

	// Read the HTTP request body
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Unable to read request body", http.StatusBadRequest)
        return
    }

	// Unmarshal the request body into the newCustomer struct
    err = json.Unmarshal(reqBody, &newCustomer)
    if err != nil {
        http.Error(w, "Invalid input format", http.StatusBadRequest)
        return
    }

    // Check if the customer ID already exists in the database
    if _, ok := db[newCustomer.ID]; ok {
        // If the ID already exists, respond with a 409 Conflict
        w.WriteHeader(http.StatusConflict)
        json.NewEncoder(w).Encode("Error: Customer with this ID already exists")
		return
    }

	// If the ID does not exist, add the new customer to the database
    db[newCustomer.ID] = newCustomer

    // Respond with a 201 Created status
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(db)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the "id" from the URL path
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Use the helper function to get the customer
	_, success := isValidId(w, idStr)
	if !success {
		return // The error is already handled in the helper function
	}

	// update customer with id: {id}
	// Decode the request body into a Customer object
	var updatedCustomer Customer
	var err = json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		// If there's an error decoding the request body, return a bad request error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if any of the required fields are empty and return bad request if so
	if updatedCustomer.Name == "" || updatedCustomer.Role == "" || updatedCustomer.Email == "" || updatedCustomer.Phone == "" {
		http.Error(w, "All fields (Name, Role, Email, Phone) are required", http.StatusBadRequest)
		return
	}

	// If the ID does not exist, add the new customer to the database
    db[updatedCustomer.ID] = updatedCustomer

    // Respond with a 201 Created status
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(db)
}

func updateCustomersBatch(w http.ResponseWriter, r *http.Request) {
    var updatedCustomers []Customer
    err := json.NewDecoder(r.Body).Decode(&updatedCustomers)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    for _, customer := range updatedCustomers {
        if _, exists := db[customer.ID]; exists {
            db[customer.ID] = customer
        } else {
            http.Error(w, "Customer not found", http.StatusNotFound)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(db)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the "id" from the URL path
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Use the helper function to get the customer
	customer, success := isValidId(w, idStr)
	if !success {
		return // The error is already handled in the helper function
	}

	delete(db, customer.ID)

	json.NewEncoder(w).Encode(db)
}

// Helper function to check customer id is valid or not, and handle errors
func isValidId(w http.ResponseWriter, idStr string) (*Customer, bool) {
	// Convert the id to uint8
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If the conversion fails, return a bad request error
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return nil, false
	}

	// Check if the customer exists in the database
	customer, exists := db[uint8(id)]
	if !exists {
		// If the customer doesn't exist, return a 404 error
		http.Error(w, "Customer not found", http.StatusNotFound)
		return nil, false
	}

	// Return the customer and a success status
	return &customer, true
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	customer1 := Customer{
		ID: 1, Name: "John Doe", 
		Role: "Subscriber", 
		Email: "john.doe@gmail.com",
		Phone: "123-456-7890",
		Contacted: true,
	}
	customer2 := Customer{
		ID: 2, Name: "Peter Pan", 
		Role: "Prospect", 
		Email: "peter.pan@gmail.com",
		Phone: "321-654-0987",
		Contacted: true,
	}
	customer3 := Customer{
		ID: 3, Name: "Mary Jane", 
		Role: "Influencer", 
		Email: "mary.jane@gmail.com",
		Phone: "111-222-3333",
		Contacted: false,
	}

	// the database includes at least three existing (i.e., "hard-coded") customers.
	db[customer1.ID] = customer1
	db[customer2.ID] = customer2
	db[customer3.ID] = customer3

	router := mux.NewRouter()

	router.HandleFunc("/", serveIndex)
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomersBatch).Methods("POST")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	
	fmt.Println("Server is starting on port 3000...")
	http.ListenAndServe(":3000", router)

}