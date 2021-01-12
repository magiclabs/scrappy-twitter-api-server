package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handleRequests() {

	// Create a new instance of a mux router ✨
	myRouter := mux.NewRouter().StrictSlash(true)

	// Home page ✨
	myRouter.HandleFunc("/", serverHomePage)

	// Return all tweets ✨
	myRouter.HandleFunc("/tweets", returnAllTweets)

	// Delete a tweet ✨
	myRouter.HandleFunc("/tweets/{id}", deleteATweet).Methods("DELETE")

	// Return a single Tweet ✨
	myRouter.HandleFunc("/tweets/{id}", returnSingleTweet)

	// Create a tweet ✨
	myRouter.HandleFunc("/tweet", createATweet).Methods("POST")

	// Save the user, if they're authenticated ✨
	myRouter.HandleFunc("/login", logIn).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	if err := http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(myRouter)); err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Tweets array will be accessible by all .go files
	Tweets = []Tweet{
		Tweet{ID: "1", Copy: "This is our first default tweet!", Author: "maricris@magic.link"},
	}

	handleRequests()
}
