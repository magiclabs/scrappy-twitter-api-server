package main

import (
	"log"
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/client"
	"github.com/magiclabs/magic-admin-go/token"
)

// Tweet is struct or data type with an Id, Copy and Author
type Tweet struct {
	ID     string `json:"ID"`
	Copy   string `json:"Copy"`
	Author string `json:"Author"`
}

// Tweets is an array of Tweet structs
var Tweets []Tweet

// User is struct with an Email
type User struct {
	Email string `json:"Email"`
}

const authBearer = "Bearer"

// Load .env file from given path
var err = godotenv.Load(".env")

// Get env variables
var magicSecretKey = os.Getenv("MAGIC_SECRET_KEY")

// Instantiate Magic ✨
var magicSDK = client.New(magicSecretKey, magic.NewDefaultClient())

// Handler for the server's homepage ✨
func serverHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to our Twitter page!")
	fmt.Println("Endpoint Hit: serverHomePage")
}

// Returns ALL tweets ✨
func returnAllTweets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllTweets")
	json.NewEncoder(w).Encode(Tweets)
}

// Returns a SINGLE tweet ✨
func returnSingleTweet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTweet")
	vars := mux.Vars(r)
	key := vars["id"]

	/*
		Loop over all of our Tweets
		If the tweet.Id equals the key we pass in
		Return the tweet encoded as JSON
	*/
	for _, tweet := range Tweets {
		if tweet.ID == key {
			json.NewEncoder(w).Encode(tweet)
		}
	}
}

// Creates a tweet ✨
func createATweet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createATweet")
	isAuthenticated := isAuthenticated(w, r)

	if isAuthenticated {

		/*
			Get the body of our POST request
			Unmarshal this into a new Tweet struct
			Append this to our Tweets array.
		*/
		reqBody, _ := ioutil.ReadAll(r.Body)
		var tweet Tweet
		json.Unmarshal(reqBody, &tweet)

		/*
			Update our global Tweets array to include
			Our new Tweet
		*/
		Tweets = append(Tweets, tweet)
		json.NewEncoder(w).Encode(tweet)

		w.Write([]byte("Yay! Tweet CREATED."))
	}
}

// Deletes a tweet ✨
func deleteATweet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteATweet")

	isAuthenticated := isAuthenticated(w, r)

	// Parse the path parameters
	vars := mux.Vars(r)

	if isAuthenticated {
		// Extract the `id` of the tweet we wish to delete
		id := vars["id"]

		// Loop through all our tweets
		for index, tweet := range Tweets {

			/*
				Checks whether or not our id path
				parameter matches one of our tweets.
			*/
			if tweet.ID == id {

				// Updates our Tweets array to remove the tweet
				Tweets = append(Tweets[:index], Tweets[index+1:]...)
			}
		}

		w.Write([]byte("Yay! Tweet has been DELETED."))
	}
}

// Ensures the access token sent by the client is valid ✨
func isAuthenticated(w http.ResponseWriter, r *http.Request) bool {

	// Check whether or not DID token exists in HTTP Header Request
	if !strings.HasPrefix(r.Header.Get("Authorization"), authBearer) {
		fmt.Fprintf(w, "Bearer token is required")
		return false
	}

	// Retrieve AUTH or DID token from HTTP Header Request
	did := r.Header.Get("Authorization")[len(authBearer)+1:]

	// Check whether or not DID is an empty string
	if did == "" {
		fmt.Fprintf(w, "DID token is required")
		return false
	}

	// What does NewToken() do?
	tk, err := token.NewToken(did)
	if err != nil {
		fmt.Fprintf(w, "Malformed DID token error: %s", err.Error())
		w.Write([]byte(err.Error()))
		return false
	}

	// Validate AUTH or DID Token
	if err := tk.Validate(); err != nil {
		fmt.Fprintf(w, "DID token failed validation: %s", err.Error())
		w.Write([]byte(err.Error()))
		return false
	}

	return true
}

// Acknowledges authenticated user upon login ✨
func logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: logIn")

	// Check whether or not DID token exists in HTTP Header Request
	if !strings.HasPrefix(r.Header.Get("Authorization"), authBearer) {
		fmt.Fprintf(w, "Bearer token is required")
		return
	}

	// Retrieve AUTH or DID token from HTTP Header Request
	did := r.Header.Get("Authorization")[len(authBearer)+1:]

	// Get body of our POST request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User

	// Unmarshal JSON data into a new User struct
	json.Unmarshal(reqBody, &user)

	/*
		Marshal User struct into JSON data to
		access key-value pair.
	*/
	json.Marshal(user)

	// Check whether or not DID is an empty string
	if did == "" {
		fmt.Fprintf(w, "DID token is required")
		return
	}

	// What does NewToken() do?
	tk, err := token.NewToken(did)
	if err != nil {
		fmt.Fprintf(w, "Malformed DID token error: %s", err.Error())
		return
	}

	// Validate AUTH or DID Token
	if err := tk.Validate(); err != nil {
		fmt.Fprintf(w, "DID token failed validation: %s", err.Error())
		return
	}

	// Get the Issuer
	// Then get the User Meta Data
	userInfo, err := magicSDK.User.GetMetadataByIssuer(tk.GetIssuer())
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	/*
		If the email sent by the client does not match
		the email saved via Magic SDK, then it is an
		unauthorized login.
	*/
	if userInfo.Email != user.Email {
		fmt.Fprintf(w, "Unauthorized user login")
		return
	}

	/*
		If you wanted, you could call your application logic to save the user's info.

		E.g.
		logic.User.add(userName, userEmail, tk.GetIssuer())
	*/

	// Instead of saving the user's info, we'll just return it
	w.Write([]byte("Yay! User was able to login / sign up. 🪄 Email: " + user.Email))
}

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
		Tweet{ID: "1", Copy: "This is our first default tweet!", Author: "Maricris"},
		Tweet{ID: "2", Copy: "This is our second default tweet!", Author: "Kona"},
	}

	handleRequests()
}
