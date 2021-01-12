package main

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
