package products

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	Product []struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		Handle      string `json:"handle"`
		Vendor      string `json:"vendor"`
		ProductType string `json:"product_type"`
		Variants    []struct {
			Id        int64  `json:"id"`
			Title     string `json:"title"`
			Sku       string `json:"sku"`
			Available bool   `json:"available"`
			Price     string `json:"price"`
		} `json:"variants"`
		Images []struct {
			Id  int64  `json:"id"`
			Src string `json:"src"`
		} `json:"images"`
	} `json:"products"`
}

// Finding and printing specfic product.
// Looks through the product and variants struct to find the specific information.
func FindProduct(s string, userId string, size string) (string, string, string, string, string, int64, bool) {
	flag := false
	list, id := getInfo(s, userId)

	// Find the product and return info.
	for _, p := range list.Product {
		if id == p.Id {
			for _, v := range p.Variants {
				if size == v.Title {
					// Make sure the product is available.
					if v.Available == false {
						continue
					} else {
						for _, i := range p.Images {
							flag = true
							return p.Title, v.Title, v.Price, p.Handle, i.Src, v.Id, flag
						}
					}
				}
			}
		}
	}
	// Return empty strings and false flag to inform user for unavailable product.
	return "", "", "", "", "", -1, flag
}

// Find the product and return the sizes that are available.
// Looks through the product and variants struct to find the specific information.
func FindSizes(s string, userId string) (string, string, string, string, string, bool) {
	flag := false
	list, id := getInfo(s, userId)
	vMap := make(map[int64]string)
	var price string

	for _, p := range list.Product {
		if id == p.Id {
			for _, v := range p.Variants {
				price = v.Price
				if v.Available == true {
					vMap[v.Id] = v.Title
				} else {
					continue
				}
			}
			sizes := sizeStrings(vMap)
			for _, i := range p.Images {
				flag = true
				return p.Title, sizes, price, p.Handle, i.Src, flag
			}
		}
	}
	// Return empty strings and false flag to inform user for unavailable product.
	return "", "", "", "", "", flag
}

// Reading json; do this before to make sure the site is a valid url.
// Converting string -> int & asking user input for product ID and size
func getInfo(s string, userId string) (Products, int64) {
	list := readURL(s)
	temp := convert(userId)
	id := validID(temp)
	return list, id
}

// To print the different sizes if that's what the user wants.
func sizeStrings(vmap map[int64]string) string {
	sizes := ""
	for _, v := range vmap {
		sizes = sizes + v + "\n"
	}
	return sizes
}

// Checks the site to see if it's a shopify site.
// Reads the json file of the site, gets the body, and parses through it.
func readURL(site string) Products {
	var item Products
	u := "https://" + site + "/products.json"
	jsonFile, err := http.Get(u)
	if err != nil {
		fmt.Println("Invalid site")
		log.Fatal(err)
	}
	defer jsonFile.Body.Close()
	body, err := ioutil.ReadAll(jsonFile.Body)

	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println("Invalid jsonFile")
		log.Fatal(err)
	}
	return item
}

// Convert string of id to int64
func convert(userid string) int64 {
	temp, err := strconv.ParseInt(userid, 0, 64)
	if err != nil {
		fmt.Println("userid is incorrect")
		log.Fatal(err)
	}
	return temp
}

// Counting the digits of product id to make sure that
// the user gives 13 digits.
func countInt(num int64) int {
	count := 0
	for num > 0 {
		num /= 10
		count++
	}
	return count
}

// Validating the product id that is given by the user.
// If not valid, it will keep asking the user until a valid id is given.
func validID(id int64) int64 {
	c := 0
	for {
		if c < 13 {
			c = countInt(id)
		} else {
			break
		}
	}
	return id
}
