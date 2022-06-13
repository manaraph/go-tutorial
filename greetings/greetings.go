package greetings

import (
	"fmt"

	"rsc.io/quote"
)

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

func RscQuote() {
	fmt.Println("Hello Manasseh")
	fmt.Println(quote.Go())
}
