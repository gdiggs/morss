package main

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/peterbourgon/diskv/v3"
)

var feedParser = gofeed.NewParser()

func checkFeed(feedUrl string, db *diskv.Diskv) {
	feed, _ := feedParser.ParseURL(feedUrl)
	log.Println("Checking updates for " + feed.Title)
	feedId := url.QueryEscape(feedUrl)

	lastUpdated := db.ReadString(feedId)

	if lastUpdated == "" || lastUpdated != feed.Updated {
		log.Println(feed.Updated + ": " + feed.Items[0].Title)

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
