package slack

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	p "json/products"
)

// This function "turns on" socketmode client, it will keep it open until we terminate the program.
func SlackCommand() {
	// IMPORTANT: Add your tokens and the channel id to the .env file or else this program will throw an error.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	auth_tok := os.Getenv("AUTH_TOKEN")
	app_tok := os.Getenv("APP_TOKEN")
	channelid := os.Getenv("CHANNEL_ID")

	// Connecting to the slack and socketmode.
	api := slack.New(auth_tok, slack.OptionDebug(true), slack.OptionAppLevelToken(app_tok))
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
	)
	fmt.Println("Slack: vbot is now running.")

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Keeps the client open to accept more request. It will close/shut if we get a log error.
	go func(c context.Context, api *slack.Client, client *socketmode.Client) {
		for {
			select {
			case <-c.Done():
				log.Println("Socketmode Shut Down.")
				return
			case event := <-client.Events:
				switch event.Type {
				case socketmode.EventTypeSlashCommand:
					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not be casted.")
						continue
					}
					// Acknowledge req
					client.Ack(*event.Request)
					err := handleSlashCommand(command, api, channelid)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}(c, api, client)
	client.Run()
}

// Handles all slash commands
func handleSlashCommand(command slack.SlashCommand, api *slack.Client, channelid string) error {
	switch command.Command {
	case "/product":
		return handleProductID(command, *api, channelid)
	case "/sizes":
		return handleSizes(command, *api, channelid)
	}
	return nil
}

// Handles accepting the product id, site, and size and plits the string into 3 diff params.
// This function finds the product by calling FindProduct from product.go
func handleProductID(command slack.SlashCommand, api slack.Client, channelid string) error {
	params := &slack.Msg{Text: command.Text}
	str := strings.Split(params.Text, " ")
	temp := deleteEmpty(str)
	site, productid, userSize := temp[0], temp[1], temp[2]

	// Finding the product in product.go
	name, size, price, handle, image, id, flag := p.FindProduct(site, productid, userSize)
	createMessage(name, size, price, image, site, handle, id, flag, api, channelid)
	return nil
}

// Handles accepting the product id and site and splits the string into 2 diff params.
// This function finds the sizes by calling FindSizes from product.go
func handleSizes(command slack.SlashCommand, api slack.Client, channelid string) error {
	params := &slack.Msg{Text: command.Text}
	str := strings.Split(params.Text, " ")
	temp := deleteEmpty(str)
	site, productid := temp[0], temp[1]

	// Finding the product in product.go
	name, size, price, handle, image, flag := p.FindSizes(site, productid)
	createMessage(name, size, price, image, site, handle, -1, flag, api, channelid)
	return nil
}

// Deletes the empty spaces in the string so our info matches our struct correctly.
func deleteEmpty(s []string) []string {
	var temp []string
	for _, str := range s {
		if str != "" {
			temp = append(temp, str)
		}
	}
	return temp
}

// This function will create an attachment in the slack channel
// that we want to post our products info to. OR Creates an "error" Attachment if our flag returns false.
func createMessage(ptitle string, vtitle string, price string, image string, s string, handle string, id int64, flag bool, api slack.Client, channelid string) {
	checkout := ""
	if flag == true {
		if id == -1 {
			checkout = "No checkout link provided."
		} else {
			checkout = "https://" + s + "/cart/" + strconv.FormatInt(id, 10) + ":1"
		}
		attachment := slack.Attachment{
			AuthorName: "https://" + s,
			Title:      ptitle,
			TitleLink:  "https://" + s + "/products/" + handle,
			Color:      "#e6e6fa",
			ThumbURL:   image,
			Fields: []slack.AttachmentField{
				{
					Title: "Size:",
					Value: vtitle + "\n",
					Short: true,
				},
				{
					Title: "Price:",
					Value: "$" + price,
					Short: true,
				},
				{
					Title: "Autocheckout Link:",
					Value: checkout,
				},
			},
			Footer: "Shivani Kumar" + " | " + time.Now().Format("01-02-2006 15:04:05 MST"),
		}
		_, timestamp, err := api.PostMessage(
			channelid,
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionAsUser(true),
		)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		log.Printf("Message successfully sent at %s\n", timestamp)
		// Error attachment if the product is not available.
	} else {
		attachment := slack.Attachment{
			AuthorName: "https://" + s,
			Title:      ptitle,
			Color:      "#ba2507",
			Fields: []slack.AttachmentField{
				{
					Title: "SORRY PRODUCT UNAVAILABLE!",
				},
			},
			Footer: "Shivani Kumar" + " | " + time.Now().Format("01-02-2006 15:04:05 MST"),
		}
		_, timestamp, err := api.PostMessage(
			channelid,
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionAsUser(true),
		)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		log.Printf("Message successfully sent at %s\n", timestamp)
	}
}
