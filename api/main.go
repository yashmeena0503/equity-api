package main

import (
	"github.com/joho/godotenv"
)

func main() {

	// Load env.
	err := godotenv.Load("../.env")
	if err != nil { println("error loading env") }

	// Initialize the App.

	println("hi")

}
