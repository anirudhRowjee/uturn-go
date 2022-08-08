package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UrlStruct struct {
	gorm.Model
	Shortcode string `json:"shortcode"`
	URL       string `json:"url"`
}

var shortcode_sample_space = []rune("abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ0123456789")
var shortcode_sample_space_length = len(shortcode_sample_space)

func generate_random_shortcode(length int) string {

	new_shortcode := make([]rune, length)
	for i := range new_shortcode {
		new_shortcode[i] = shortcode_sample_space[rand.Int63n(int64(shortcode_sample_space_length))]
	}
	return string(new_shortcode)
}

func main() {

	// seed the random number generator
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Random text >> ", generate_random_shortcode(10))

	// Test static data
	var store_static = []UrlStruct{
		{Shortcode: "abc", URL: "https://www.google.com"},
		{Shortcode: "def", URL: "https://www.amazon.com"},
		{Shortcode: "ghi", URL: "https://www.twitter.com"},
		{Shortcode: "lmno", URL: "https://www.facebook.com"},
	}
	fmt.Println(len(store_static))

	// Adding database connectivity
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Could not open Database!")
	}

	// Automigrate the schema
	db.AutoMigrate(&UrlStruct{})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/urls", func(c *gin.Context) {

		// get all the URLs in JSON.
		var results []UrlStruct
		db.Find(&results)
		c.JSON(200, results)

	})

	r.POST("/urls", func(c *gin.Context) {

		var currentUrl UrlStruct

		// parse request body to get the proposed URL Object
		if err := c.BindJSON(&currentUrl); err != nil {
			fmt.Println("Error ! Could not read JSON Formatting!")
			c.JSON(400, gin.H{
				"Error": "Improperly shaped JSON request body",
			})
		}

		if currentUrl.Shortcode == "" {
			currentUrl.Shortcode = generate_random_shortcode(5)
		}

		// insert this into the database
		result := db.Create(&currentUrl)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error":  "Could not insert into database",
				"Reason": result.Error,
			})
		}

		// send a success confirmation
		c.JSON(200, currentUrl)
	})

	r.GET("/:shortcode", func(c *gin.Context) {

		// get the shortcode from the URL Body
		shortcode, success := c.Params.Get("shortcode")
		if success == false {
			c.JSON(404, gin.H{
				"Error": "Could not parse URL",
			})
		}

		// fetch the shortcode from the database
		current_url := UrlStruct{Shortcode: shortcode}
		db.Where("shortcode = ?", shortcode).First(&current_url)

		fmt.Println("Current URL Fetched: ", current_url)

		// send a redirect to there
		c.Redirect(301, current_url.URL)
	})

	r.Run(":8000")
}
