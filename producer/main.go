package main

import (
	"log"
	"net/http"
)

type Order struct {
	CustomerName string `json: "customer_name"`
	CoffeeType   string `json: "coffee_type"`
}

func main() {
	//	fmt.Println("Hello, World!")
	http.HandleFunc("/order", placeOrder)
	log.Fatal(http.ListenAndServe(":8400", nil))
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
