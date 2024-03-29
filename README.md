# Scrappy Twitter API
Scrappy Twitter API is a Go-backend project that is secured by the Magic SDK for Go. 

# Scrappy Twitter API (SERVER)
This Go server is where all of the Scrappy Twitter API requests are handled. Once the user has generated a [Decentralised ID Token (DIDT)](https://docs.magic.link/decentralized-id) from the [client side](https://github.com/magiclabs/scrappy-twitter-api-client), they can pass it into their Request Header as a Bearer token to hit protected endpoints.

# API Routes
- POST a tweet (protected): http://localhost:8080/tweet 
- GET all tweets (unprotected): http://localhost:8080/tweets 
- GET a single tweet (unprotected): http://localhost:8080/tweet/1 
- DELETE a tweet (protected): http://localhost:8080/tweet/2

# Noteworthy Dependencies
- [gorilla/handlers](https://github.com/gorilla/handlers): Lets us enable CORS.
- [gorilla/mux](https://github.com/gorilla/mux): Lets us build a powerful HTTP router and URL matcher.
- [magic-admin-go/client](https://docs.magic.link/admin-sdk/go/get-started#creating-an-sdk-client-instance): Lets us instantiate the Magic SDK for Go.
- [magic-admin-go/token](https://docs.magic.link/admin-sdk/go/get-started#creating-a-token-instance): Lets us create a Token instance.

# Quickstart
## Magic Setup
1. Sign up for an account on [Magic](https://magic.link/).
2. Create an app.
3. Copy your app's Live Secret Key (you'll need it soon).

## Server Setup
1. `git clone https://github.com/seemcat/scrappy-twitter-api-server.git`
2. `cd scrappy-twitter-api-server`
3. `mv .env.example .env`
4. Paste the Live Secret Key you just copied as the value for `MAGIC_SECRET_KEY` in .env:
    ```
    MAGIC_SECRET_KEY=sk_XXXXXXXXXX
    ```
4. Run all .go files with `go run .`

## Test with Postman
1. Import the DEV version of the Scrappy Twitter API Postman Collection:
    [![Run in Postman](https://run.pstmn.io/button.svg)](https://god.postman.co/run-collection/1aa913713995cb16bb70)
2. Generate a DID token on the Client side. 
   
   Click [here](https://github.com/magiclabs/scrappy-twitter-api-client) to spin up your own local client and generate the DID token there.
   
3. Pass the DID token as a Bearer token into the Postman Collection’s HTTP Authorization request header and click **save**.
4. Send your requests to the Scrappy Twitter API! 🎉
