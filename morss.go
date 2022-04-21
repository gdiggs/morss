package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/mmcdole/gofeed"
	"github.com/peterbourgon/diskv/v3"
)

var feedParser = gofeed.NewParser()

func pingSlack(title string, message string, link string) {
	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")

	attachment := slack.Attachment{}
	attachment.AddAction(slack.Action{Type: "button", Text: "Link", Url: link, Style: "primary"})

	messageText := fmt.Sprintf("*%s*: %s", title, message)
	payload := slack.Payload{
		Text:        messageText,
		Username:    "morss",
		IconEmoji:   ":radio:",
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		log.Fatal(err)
		os.Exit(1)
	}
}

func checkFeed(feedUrl string, db *diskv.Diskv) {
	feed, _ := feedParser.ParseURL(feedUrl)
	log.Println("Checking updates for " + feed.Title)
	feedId := url.QueryEscape(feedUrl)

	lastUpdated := db.ReadString(feedId)

	if lastUpdated == "" || lastUpdated != feed.Updated {
		item := feed.Items[0]
		log.Println(feed.Updated + ": " + item.Title)
		pingSlack(feed.Title, item.Title, item.Link)

		err := db.Write(feedId, []byte(feed.Updated))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		log.Println("No updates")
	}
}

func main() {
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     os.Getenv("DATASTORE"),
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	urls := os.Getenv("FEED_URLS")
	if urls == "" {
		log.Fatal("Error parsing Feed URLs")
		os.Exit(1)
	}
	feeds := strings.Split(urls, ",")

	for {
		for _, url := range feeds {
			checkFeed(url, db)
		}
		time.Sleep(60 * time.Second)
	}
}
