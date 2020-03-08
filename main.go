package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
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
		ocagent.WithReconnectionPeriod(10*time.Second),
		ocagent.WithAddress(ocagentHost),
		ocagent.WithServiceName("welcomer"))

	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	r := gin.Default()

	r.GET("/welcome", func(c *gin.Context) {
		//_, span := trace.StartSpan(c, "/welcome")
		// http_server_route=/welcome tag is set
		ochttp.SetRoute(c.Request.Context(), "/welcome")

		//defer span.End()
		fmt.Println(c.Request.Header)
		fmt.Println(c.Request.Host)

		welcomeHandler(c)
	})
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	http.ListenAndServe( // nolint: errcheck
		"0.0.0.0:8080",
		&ochttp.Handler{
			Handler: r,
			GetStartOptions: func(r *http.Request) trace.StartOptions {
				startOptions := trace.StartOptions{}

				if r.URL.Path == "/metrics" {
					startOptions.Sampler = trace.NeverSample()
				}

				return startOptions
			},
		}, )
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
	context := c.Request.Context()
	span := trace.FromContext(context)
	defer span.End()
	span.Annotate([]trace.Attribute{trace.StringAttribute("annotated", "welcomervalue")}, "welcomervalue-->guesttracker annotation check")
	span.AddAttributes(trace.StringAttribute("span-add-attribute", "welcomervalue"))
	time.Sleep(time.Millisecond * 125)

	r, _ := http.NewRequest("POST", "http://"+guestrackerhost+"/track-guest", bytes.NewBuffer(reqBody))
	clientTrace := ochttp.NewSpanAnnotatingClientTrace(r, span)
	context = httptrace.WithClientTrace(context, clientTrace)
	r = r.WithContext(context)

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
