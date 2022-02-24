package discord

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	p "json/products"
)

// Turns the connection to discord on.
func DiscordCommand() {
	// IMPORTANT: Add your oauth token to the .env file or else this program will throw an error.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	token := os.Getenv("TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(reply)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("vbot is now running.")

	// If you want to only run the discord bot alone. Uncomment the following.
	/*terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-terminate

	dg.Close() */
}

// Read in user input and spilt the string so we can pass to findProduct
// If our command is !findproduct or pass to findSize if our command is !findSizes.
func reply(sess *discordgo.Session, msg *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself.
	if msg.Author.ID == sess.State.User.ID {
		return
	}

	if strings.Contains(msg.Content, "!findproduct") {
		str := strings.Split(msg.Content, " ")
		temp := deleteEmpty(str)
		site, productid, userSize := temp[1], temp[2], temp[3]

		name, size, price, image, site, handle, id, flag := findProduct(site, productid, userSize)
		postMessage(name, size, price, image, site, handle, id, flag, sess, msg)
	}

	if strings.Contains(msg.Content, "!findsize") {
		str := strings.Split(msg.Content, " ")
		temp := deleteEmpty(str)
		site, productid := temp[1], temp[2]

		name, size, price, image, site, handle, id, flag := findSize(site, productid)
		postMessage(name, size, price, image, site, handle, id, flag, sess, msg)
	}
}

// This function finds the product by calling FindProduct from product.go
func findProduct(site string, productid string, userSize string) (string, string, string, string, string, string, int64, bool) {
	name, size, price, handle, image, id, flag := p.FindProduct(site, productid, userSize)
	return name, size, price, image, site, handle, id, flag
}

// This function finds the sizes by calling FindSizes from product.go
func findSize(site string, productid string) (string, string, string, string, string, string, int64, bool) {
	name, size, price, handle, image, flag := p.FindSizes(site, productid)
	return name, size, price, image, site, handle, -1, flag
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

// Create the attachment that we will be posting to discord and sends the message to the channel.
// that contains our products information if our flag returns true. OR Creates an "error" Attachment if our flag returns false.
func postMessage(name string, size string, price string, image string, site string, handle string, id int64, flag bool, sess *discordgo.Session, msg *discordgo.MessageCreate) {
	checkout := ""
	if flag == true {
		if id == -1 {
			checkout = "No checkout link provided."
		} else {
			checkout = "https://" + site + "/cart/" + strconv.FormatInt(id, 10) + ":1"
		}
		attachment := &discordgo.MessageEmbed{
			Author: (&discordgo.MessageEmbedAuthor{
				Name: ("https://" + site),
			}),
			Color: 0xe6e6fa,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Price",
					Value:  "$" + price,
					Inline: true,
				},
				{
					Name:   "Size(s)",
					Value:  size,
					Inline: true,
				},
				{
					Name:  "Auto-checkout Link",
					Value: checkout,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: image,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     name,
			URL:       ("https://" + site + "/products/" + handle),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Shivani Kumar",
			},
		}
		_, err := sess.ChannelMessageSendEmbed(msg.ChannelID, attachment)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// Give an error in case the product the user wanted was not found.
		attachment := &discordgo.MessageEmbed{
			Author: (&discordgo.MessageEmbedAuthor{
				Name: ("https://" + site),
			}),
			Color: 0xba2507,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "PRODUCT UNAVAILABLE.",
					Value: "The product you are looking for is currently unavailable.",
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: image,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     name,

			Footer: &discordgo.MessageEmbedFooter{
				Text: "Shivani Kumar",
			},
		}
		_, err := sess.ChannelMessageSendEmbed(msg.ChannelID, attachment)
		if err != nil {
			fmt.Println(err)
		}
	}
}
