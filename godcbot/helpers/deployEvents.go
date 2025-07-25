package helpers

import (
	"log"

	"github.com/Sush1sui/dcbot/events"
	"github.com/bwmarrin/discordgo"
)


var EventHandlers = map[string]interface{}{
	"MessageCreate": events.OnPing,
	// Add more event handlers here, e.g.:
	// Go doesn't support dynamic runtime imports
	// You have to manually add each event handler
}

func DeployEvents(sess *discordgo.Session) {
	for _, handler := range EventHandlers {
		sess.AddHandler(handler)
	}
	log.Println("Event handlers deployed successfully.")
}