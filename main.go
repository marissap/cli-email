package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readInput(p string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	s, _ := reader.ReadString('\n')
	return s
}

func main() {

	colorCyan := "\033[36m"
	colourRed := "\033[31m"
	var email, password, smtp, pin, useExisting string

	// check for existing cache and load if true
	cache := &Cache{}
	c, err := cache.LoadCacheFromFile()
	if err != nil {
		fmt.Println(string(colorCyan), "\nYou don't have any saved email accounts, let's start from scratch!")
		useExisting = "N"
	}

	if useExisting == "N" {
		email = strings.TrimSuffix(readInput("\nPlease enter your email: "), "\n")
		password = strings.TrimSuffix(readInput("\nPlease enter your password (Don't worry, I won't look): "), "\n")
		pin = strings.TrimSuffix(readInput("\nPlease enter a alphanumeric pin of length 4 (This allows us to encrypt and decrypt your login): "), "\n")
	} else {
		pin = strings.TrimSuffix(readInput("\nPlease enter your pin: "), "\n")
		email = c.items[pin].email
		smtp = c.items[pin].smtp
		// useExisting = strings.ToUpper(strings.TrimSuffix(readInput("Would you like to use existing email: "+email+"? [y/n]"), "\n"))
		password, err = Decrypt(pin, c.items[pin].pwd)
		if err != nil {
			fmt.Println(string(colourRed), "\nUnable to decrypt password. Try again later!")
			os.Exit(3)
		}
	}

	fmt.Println(string(colorCyan), "\nOkay! Let's start writing!📝")
	recieverName := strings.TrimSuffix(readInput("\nWhat is the name of the person you are emailing?"), "\n")
	recieverEmail := strings.TrimSuffix(readInput("\nWhat is their email address?"), "\n")
	messageToSend := strings.TrimSuffix(readInput("\nSounds good! Start writing!"), "\n")
	fmt.Println(string(colorCyan), "\nYou sound like Shakespeare! Attempting to send email now... 📧")

	// send email here
	err = Send(email, recieverEmail, smtp, password, recieverName, messageToSend)
	if err != nil {
		fmt.Println(string(colourRed), "\nError sending email. Sorry, try again later!")
	}

}
