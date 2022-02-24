package main

import (
	d "json/discord"
	s "json/slack"
)

func main() {
	d.DiscordCommand()
	s.SlackCommand()
}
