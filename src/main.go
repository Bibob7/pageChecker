package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	app := &cli.App{
		Name:  "pcheck",
		Usage: "Page Checker to observe changes on web pages.",
		Commands: cli.Commands{
			{
				Name: "observe",
				Action: func(c *cli.Context) error {
					uri := c.String("uri")
					webHookUri := c.String("webHookuri")
					listenForChangesOn(uri, webHookUri)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "uri",
						Aliases:  []string{"u"},
						Usage:    "Uri of the content you want to observe",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "webHookuri",
						Aliases:  []string{"wu"},
						Usage:    "Weebhook Url to push notification to.",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listenForChangesOn(url string, weebHookUrl string) {
	done := make(chan bool)

	log.Printf("Start observing page: %s.", url)

	go func() {
		for {
			if !checkPageForChanges(url, weebHookUrl) {
				time.Sleep(10 * time.Second)
				continue
			}
			done <- true
		}
	}()

	<-done
	log.Printf("Listening for changes on %s stopped.", url)
}

func checkPageForChanges(url string, weebhookUrl string) bool {
	log.Println("Request sent")
	response, err := http.Get(url)
	if err != nil {
		log.Printf("No response from url with error: %s", err)
		return false
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Not possible to read body with error: %s", err)
		return false
	}
	body := string(bodyBytes)
	if strings.Contains(body, "Derzeit nicht verfügbar.") {
		log.Println("No content changes yet.")
		return false
	}
	log.Println("Send slack notification and finish.")
	_ = SendSlackNotification(
		weebhookUrl,
		fmt.Sprintf("Page under %s has changed!", url),
		"")
	return true
}
