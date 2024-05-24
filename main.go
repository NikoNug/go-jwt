package main

import (
	"golang-jwt-practice/controllers/authcontroller"
	"golang-jwt-practice/controllers/productcontroller"
	"golang-jwt-practice/middlewares"
	"golang-jwt-practice/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	models.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
