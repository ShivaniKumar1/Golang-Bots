# Slack and Discord Golang Bots
The following bots were written in Golang and run from the same main.go file. They run simultaneously since both have features of connecting to the servers without interfering with each other.

## Purpose
To be able to fetch json information from all shopify sites for specifically asked products. Given a site, product id, and a size, it will return an auto checkout link within a couple seconds for the user to checkout. Given a site and product id, it will return different sizes that are available for that specific product. This is helpful when the user isn’t specifically looking for a product in one size but would like to know what sizes are available.

## Hosting
The bots are hosted and deployed on Docker and as a droplet on DigitalOcean.

## How it works
When given a site, productId, and size, all three are validated to confirm correct working urls. Then, the bots browse through the json file of that site, find a productId that matches with the user id input, confirm the product is available as well as the size is correct, and use the name, price, size, and id to generate an auto checkout link. This information gets used to create the message to be sent to the discord server or the slack server. The same thing happens with finding the sizes; however, no checkout link is provided, only a link to preview the product. If a product is not found, a message will display for the user letting them know that the product they are looking for was not found. Slack handles all commands as slash commands while discord looks for the specific command in the user’s input.
