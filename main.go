package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rest/controllers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{username}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/users/{username}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{username}", controllers.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
