# Scrappy Twitter API
Scrappy Twitter API is a Go-backend project that is secured by the Magic SDK for Go. 

# Scrappy Twitter API (SERVER)
This Go server is where all of the Scrappy Twitter API requests are handled. Once the user has generated an access token from the [client side](https://github.com/seemcat/scrappy-twitter-api-client), they can pass it into their Request Header as a Bearer token to hit protected endpoints.

# API Routes
- POST a tweet (protected): http://localhost:8080/tweet 
- GET all tweets (unprotected): http://localhost:8080/tweets 
- GET a single tweet (unprotected): http://localhost:8080/tweets/1 
- DELETE a tweet (protected): http://localhost:8080/tweets/2

# Noteworthy Packages
- [gorilla/handlers](https://github.com/gorilla/handlers): Lets us enable CORS.
- [gorilla/mux](https://github.com/gorilla/mux): Lets us build a powerful HTTP router and URL matcher.
- [magic-admin-go/client](https://docs.magic.link/admin-sdk/go/get-started#creating-an-sdk-client-instance): Lets us instantiate the Magic SDK for Go.
- [magic-admin-go/token](https://docs.magic.link/admin-sdk/go/get-started#creating-a-token-instance): Lets us create a Token instance.

# Quickstart
## Magic Setup
1. Sign up for an account on [Magic](https://magic.link/).
2. Create an app.
3. Copy your app's Test Secret Key (you'll need it soon).

## Server Setup
1. `git clone https://github.com/seemcat/scrappy-twitter-api-server.git`
2. `cd scrappy-twitter-api-server`
3. `mv .env.example .env`
4. Paste the Test Secret Key you just copied as the value for `MAGIC_TEST_SECRET_KEY` in .env:
    ```
    MAGIC_TEST_SECRET_KEY=sk_test_XXXXXXXXXX
    ```
4. Run all .go files with `go run .`

## Test with Postman
1. Import the DEV version of the Scrappy Twitter API Postman Collection:
    [![Run in Postman](https://run.pstmn.io/button.svg)](https://god.postman.co/run-collection/1aa913713995cb16bb70)
2. Generate an access token on the Client side. 
    
    **Note**: You have two options to do this. You could either click [here](https://github.com/seemcat/scrappy-twitter-api-client) to spin up your own local client and generate the access token there. 
    
    OR you could visit the **Live** client side [here](https://scrappy-twitter-api-client.vercel.app/) and immediately generate your access token there.
3. Pass the access token as a Bearer token into the Postman Collection’s HTTP Authorization request header.
4. Send your requests to the Scrappy Twitter API! 🎉