package main

import (
	d "json/discord"
	s "json/slack"
)

func main() {
	// Running Discord and Slack bots
	d.DiscordCommand()
	s.SlackCommand()
}
