package main

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/client"
	"github.com/magiclabs/magic-admin-go/token"
)

const authBearer = "Bearer"

// Load .env file from given path
var err = godotenv.Load(".env")

// Get env variables
var magicSecretKey = os.Getenv("MAGIC_TEST_SECRET_KEY")

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

	author := getAuthor(w, r)

	// Ensure the author of the tweet is authenticated
	if len(author) > 0 {

		/*
			Get the body of our POST request
			Unmarshal this into a new Tweet struct
			Add the authenticated author to the tweet
		*/
		reqBody, _ := ioutil.ReadAll(r.Body)
		var tweet Tweet
		json.Unmarshal(reqBody, &tweet)
		tweet.Author = author

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

	author := getAuthor(w, r)

	// Parse the path parameters
	vars := mux.Vars(r)

	// Ensure the author of the tweet is authenticated
	if len(author) > 0 {
		// Extract the `id` of the tweet we wish to delete
		id := vars["id"]

		// Loop through all our tweets
		for index, tweet := range Tweets {

			/*
				Checks whether or not our id and author path
				parameter matches one of our tweets.
			*/
			if (tweet.ID == id) && (tweet.Author == author) {

				// Updates our Tweets array to remove the tweet
				Tweets = append(Tweets[:index], Tweets[index+1:]...)
				w.Write([]byte("Yay! Tweet has been DELETED."))
				return
			}
		}
	}

	w.Write([]byte("Ooh. You can't delete someone else's tweet."))
}

/*
Ensures the access token sent by the client is valid
Returns the tweet author's email address ✨
*/
func getAuthor(w http.ResponseWriter, r *http.Request) string {

	// Check whether or not access token exists in HTTP Header Request
	if !strings.HasPrefix(r.Header.Get("Authorization"), authBearer) {
		fmt.Fprintf(w, "Bearer token is required")
		return ""
	}

	// Retrieve access token from HTTP Header Request
	accessToken := r.Header.Get("Authorization")[len(authBearer)+1:]

	// Check whether or not access token is an empty string
	if accessToken == "" {
		fmt.Fprintf(w, "Access token is required")
		return ""
	}

	// Create a Token instance to interact with the DID token
	tk, err := token.NewToken(accessToken)
	if err != nil {
		fmt.Fprintf(w, "Malformed access token error: %s", err.Error())
		w.Write([]byte(err.Error()))
		return ""
	}

	// Validate the Token instance before using it
	if err := tk.Validate(); err != nil {
		fmt.Fprintf(w, "DID token failed validation: %s", err.Error())
		return ""
	}

	// Get the the user's information
	userInfo, err := magicSDK.User.GetMetadataByIssuer(tk.GetIssuer())
	if err != nil {
		fmt.Fprintf(w, "Error when calling GetMetadataByIssuer: %s", err.Error())
		return ""
	}

	return userInfo.Email
}

// Acknowledges authenticated user upon login ✨
func logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: logIn")

	// Check whether or not access token exists in HTTP Header Request
	if !strings.HasPrefix(r.Header.Get("Authorization"), authBearer) {
		fmt.Fprintf(w, "Bearer token is required")
		return
	}

	// Retrieve access token from HTTP Header Request
	accessToken := r.Header.Get("Authorization")[len(authBearer)+1:]

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

	// Check whether or not access token is an empty string
	if accessToken == "" {
		fmt.Fprintf(w, "access token is required")
		return
	}

	// Create a Token instance to validate access token
	tk, err := token.NewToken(accessToken)
	if err != nil {
		fmt.Fprintf(w, "Malformed access token error: %s", err.Error())
		return
	}

	// Validate access token
	if err := tk.Validate(); err != nil {
		fmt.Fprintf(w, "access token failed validation: %s", err.Error())
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
