package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENTRA")
	json.NewEncoder(w).Encode([]byte(`{"funciona":"Ok"}`))
}

func main() {
	r := mux.NewRouter()

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "public-api-key", "X-Idempotency-Key", "private-secret-key"})
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	r.HandleFunc("/", test)
	r.HandleFunc("/v1/payments", yunoPay).Methods("POST")

	corsHandler := handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r)

	err := http.ListenAndServe(":8081", corsHandler)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func yunoPay(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	publicKey := r.Header.Get("public-api-key")
	privateKey := r.Header.Get("private-secret-key")

	url := "https://api-sandbox.y.uno/v1/payments"
	payloadByte, _ := io.ReadAll(r.Body)
	fmt.Println(string(payloadByte))
	payload := strings.NewReader(string(payloadByte))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("charset", "utf-8")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("public-api-key", publicKey)
	req.Header.Add("private-secret-key", privateKey)
	req.Header.Add("X-Idempotency-Key", idempotencyKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	response, _ := io.ReadAll(res.Body)

	w.Write(response)
}
