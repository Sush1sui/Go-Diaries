package main

import (
	"booking-app/helper"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Package Level Variables aka Global Variables
var conferenceName = "Go Conference"
const conferenceTickets int = 50
var remainingTickets uint = 50
// var bookings = [50]string{"Nana", "Kofi", "Ama", "Kwame", "Akosua"}
// var bookings [50]string

// Slice (Dynamic array)
// var bookings []string

// Array of maps
// for easier understanding coming from a JavaScript background
// maps are like objects in JS
// maps are key-value pairs
// var bookings =  make([]map[string]string, 0)
var bookings = make([]UserData, 0)

// custome type
type UserData struct {
	firstName string
	lastName  string
	email     string
	ticketCount uint
}

var wg = sync.WaitGroup{}

func main() {
	greetUsers()

	// for remainingTickets > 0 && len(bookings) < conferenceTickets {

		firstName, lastName, email, userTickets := getUserInput()

		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if !isValidName || !isValidEmail || !isValidTicketNumber {
			if !isValidName {
				fmt.Println("First name or last name is too short. Please try again.")
			}
			if !isValidEmail {
				fmt.Println("Email address is not valid. Please try again.")
			}
			if !isValidTicketNumber {
				fmt.Printf("Number of tickets you entered is invalid. Please enter a number between 1 and %v.\n", remainingTickets)
			}
			// continue
		}

		bookTicket(firstName, lastName, email, userTickets)
		// fmt.Printf("List of bookings is: %v\n", bookings)
		wg.Add(1)
		go sendTicket(firstName, lastName, email, userTickets)

		firstNames := getFirstNames()
		fmt.Printf("The first names of bookings are: %v\n", firstNames)

		noTicketsRemaining := remainingTickets == 0
		if noTicketsRemaining {
			fmt.Println("Our conference is booked out. Come back next year!")
			// break
		}
	// }
	wg.Wait() // Wait for all goroutines to finish
}


func greetUsers() {
	// fmt.Printf("conferenceTickets is %T, remainingTickets is %T, conferenceName is %T\n", confTickets, remTickets, confName)
	fmt.Printf("Welcome to %v booking application!\n", conferenceName)
	fmt.Printf("We have a total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend!")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your first name:")
	firstName, _ = reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)

	fmt.Println("Enter your last name:")
	lastName, _ = reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Println("Enter your email:")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Println("Enter number of tickets:")
	userTicketsInput, _ := reader.ReadString('\n')
	userTicketsInput = strings.TrimSpace(userTicketsInput)
	
	// Convert string to uint
	userTicketsInt, err := strconv.ParseUint(userTicketsInput, 10, 64)
	if err != nil {
		userTickets = 0 // Set to 0 if invalid, validation will catch this
	} else {
		userTickets = uint(userTicketsInt)
	}

	return firstName, lastName, email, userTickets
}

func bookTicket(firstName, lastName, email string, userTickets uint) {
	remainingTickets -= userTickets

	// create a map for a user
	var userData = UserData {
		firstName: firstName,
		lastName: lastName,
		email: email,
		ticketCount: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("Thank you %v %v for booking %v tickets. A confirmation will be sent to your email %v.\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v Tickets remaining.\n", remainingTickets)
}

func sendTicket(firstName, lastName, email string, userTickets uint) {
	// Random delay between 15-20 seconds
	randomSeconds := rand.Intn(6) + 15
	time.Sleep(time.Duration(randomSeconds) * time.Second)

	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("########################")
	fmt.Printf("Sending ticket:\n %v\nto email address %v\n", ticket, email)
	fmt.Println("########################")
	fmt.Println("Time taken to send ticket:", randomSeconds, "seconds")
	wg.Done() // Signal that this goroutine is done
}