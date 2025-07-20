package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/dcbot/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Create a new Discord session using the provided bot token
	sess, err := discordgo.New("Bot " + os.Getenv("bot_token"))
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}



	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	defer sess.Close()

	// Register the ping command handler
	helpers.DeployCommands(sess)

	// Register the event handlers
	helpers.DeployEvents(sess)

	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}