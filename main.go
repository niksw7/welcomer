package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
)

var (
	guestrackerhost = os.Getenv("GUEST_TRACKER_HOST")
)

func main() {
	if guestrackerhost == "" {
		guestrackerhost = "localhost"
	}
	fmt.Println("GUEST_TRACKER_HOST =", guestrackerhost)

	ocagentHost := "oc-collector.tracing:55678"
	oce, _ := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(1*time.Second),
		ocagent.WithAddress(ocagentHost),
		ocagent.WithServiceName("welcomer"))

	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})


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
	client := &http.Client{Transport: &ochttp.Transport{}}
	ctxnew, span := trace.StartSpan(c.Request.Context(), "child")
	fmt.Println("-------------------")
	fmt.Println(ctxnew)
	defer span.End()
	span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
	span.AddAttributes(trace.StringAttribute("hello", "world"))
	time.Sleep(time.Millisecond * 125)
	
	r, _ := http.NewRequest("POST", "http://"+guestrackerhost+"/track-guest", bytes.NewBuffer(reqBody))

	r = r.WithContext(c.Request.Context())

	// resp, err := http.Post("http://"+guestrackerhost+"/track-guest",
	// 	"application/json", bytes.NewBuffer(reqBody))

	resp, err := client.Do(r)

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
