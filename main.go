package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	guestrackerhost = os.Getenv("GUEST_TRACKER_HOST")
)

func main() {
	if guestrackerhost == "" {
		guestrackerhost = "localhost"
	}

	fmt.Println("GUEST_TRACKER_HOST =",guestrackerhost)
	
	r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) {
		fmt.Println(c.Request.Header)
		fmt.Println(c.Request.Host)
		welcomeHandler(c)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func welcomeHandler(c *gin.Context) {
	//Send post request to another service
	guesttracker(c)
	c.JSON(200, gin.H{
		"message": "Hello Folks .. You are welcome(Shhh... and also tracked by guesttracker)!!",
	})
}

func guesttracker(c *gin.Context) {
	reqBody, err := json.Marshal(map[string]string{
		"username": "Bruce Wayne",
		"email":    "batman@loreans.com",
	})
	if err != nil {
		print(err)
	}
	resp, err := http.Post("http://"+guestrackerhost+":8081/track-guest",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))

}
