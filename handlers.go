package main

import (
	"net/http"

	"context"
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

type httpHandlerFunc func(http.ResponseWriter, *http.Request)
type key string

const userInfoKey key = "userInfo"
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

	// Get the authenticated author's info from context values
	userInfo := r.Context().Value(userInfoKey)
	userInfoMap := userInfo.(*magic.UserInfo)

	/*
		Get the body of our POST request
		Unmarshal this into a new Tweet struct
		Add the authenticated author to the tweet
	*/
	reqBody, _ := ioutil.ReadAll(r.Body)
	var tweet Tweet
	json.Unmarshal(reqBody, &tweet)
	tweet.Author = userInfoMap.Email

	/*
		Update our global Tweets array to include
		Our new Tweet
	*/
	Tweets = append(Tweets, tweet)
	json.NewEncoder(w).Encode(tweet)

	fmt.Println("Yay! Tweet CREATED.")
}

// Deletes a tweet ✨
func deleteATweet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteATweet")

	// Get the authenticated author's info from context values
	userInfo := r.Context().Value(userInfoKey)
	userInfoMap := userInfo.(*magic.UserInfo)

	// Parse the path parameters
	vars := mux.Vars(r)
	// Extract the `id` of the tweet we wish to delete
	id := vars["id"]

	// Loop through all our tweets
	for index, tweet := range Tweets {

		/*
			Checks whether or not our id and author path
			parameter matches one of our tweets.
		*/
		if (tweet.ID == id) && (tweet.Author == userInfoMap.Email) {

			// Updates our Tweets array to remove the tweet
			Tweets = append(Tweets[:index], Tweets[index+1:]...)
			w.Write([]byte("Yay! Tweet has been DELETED."))
			return
		}
	}

	w.Write([]byte("Ooh. You can't delete someone else's tweet."))
}

/*
Ensures the Decentralised ID Token (DIDT) sent by the client is valid
Saves the author's user info in context values ✨
*/
func checkBearerToken(next httpHandlerFunc) httpHandlerFunc {
	fmt.Println("Middleware Hit: checkBearerToken")
	return func(res http.ResponseWriter, req *http.Request) {

		// Check whether or not DIDT exists in HTTP Header Request
		if !strings.HasPrefix(req.Header.Get("Authorization"), authBearer) {
			fmt.Fprintf(res, "Bearer token is required")
			return
		}

		// Retrieve DIDT token from HTTP Header Request
		didToken := req.Header.Get("Authorization")[len(authBearer)+1:]

		// Create a Token instance to interact with the DID token
		tk, err := token.NewToken(didToken)
		if err != nil {
			fmt.Fprintf(res, "Malformed DID token error: %s", err.Error())
			res.Write([]byte(err.Error()))
			return
		}

		// Validate the Token instance before using it
		if err := tk.Validate(); err != nil {
			fmt.Fprintf(res, "DID token failed validation: %s", err.Error())
			return
		}

		// Get the the user's information
		userInfo, err := magicSDK.User.GetMetadataByIssuer(tk.GetIssuer())
		if err != nil {
			fmt.Fprintf(res, "Error when calling GetMetadataByIssuer: %s", err.Error())
			return
		}

		// Use context values to store user's info
		ctx := context.WithValue(req.Context(), userInfoKey, userInfo)
		req = req.WithContext(ctx)
		next(res, req)
	}
}

// Acknowledges authenticated user upon login ✨
func logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: logIn")

	// Get the authenticated author's info from context values
	userInfo := r.Context().Value(userInfoKey)
	userInfoMap := userInfo.(*magic.UserInfo)

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

	/*
		If the email sent by the client does not match
		the email saved via Magic SDK, then it is an
		unauthorized login.
	*/
	if userInfoMap.Email != user.Email {
		fmt.Fprintf(w, "Unauthorized user login")
		return
	}

	/*
		If you wanted, you could call your application logic to save the user's info.

		E.g.
		logic.User.add(userInfoMap.Email, userInfoMap.Issuer, userInfo.PublicAddress)
	*/

	// Instead of saving the user's info, we'll just return it
	w.Write([]byte("Yay! User was able to login / sign up. 🪄 Email: " + user.Email))
}
