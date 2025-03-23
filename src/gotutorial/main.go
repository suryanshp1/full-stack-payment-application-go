package main

import ("fmt"
		"net/http"
		"log"
		"encoding/json"
		"bytes"
		"io"
		"github.com/stripe/stripe-go/v74"
		"github.com/stripe/stripe-go/v74/paymentintent"
);

func main() {
	stripe.Key = "xbwxhbwhxbhwbxwxhw"
	http.HandleFunc("/create-payment-intent/", handleCreatePaymentIntent)
	http.HandleFunc("/health/", handleHealthCheck)
	fmt.Println("Server is running on port 4242...")
	var err error = http.ListenAndServe("localhost:4242", nil)

	if err != nil {
		log.Fatal("Something went wrong: ", err);
		return
	}

}

func handleCreatePaymentIntent(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		ProductID string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
		Address1 string `json:"address_1"`
		Address2 string `json:"address_2"`
		City string `json:"city"`
		State string `json:"state"`
		Zip string `json:"zip"`
		Country string `json:"country"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)

	if err != nil {
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(calculateOrderAmount(req.ProductID)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params)

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// fmt.Println(paymentIntent.ClientSecret)

	var response struct {
		ClientSecret string `json:"client_secret"`
	}

	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(writer, &buf)

	if err != nil {
		log.Fatal("Something went wrong: ", err)
		return
	}
	
}

func handleHealthCheck(writer http.ResponseWriter, request *http.Request) {
	var responseString string = "Server is up and running!"
	// var response []byte = []byte(responseString)
	response := []byte(responseString)
	_, err := writer.Write(response)

	if err != nil {
		log.Fatal("Something went wrong: ", err)
		return
	}
}


func calculateOrderAmount(productID string) int64 {
	// Replace this constant with a call to your database to retrieve the price for the product
	switch productID {
	case "1":
		return 109
	case "2":
		return 209
	case "3":
		return 309
	default:
		return 0
	}
}