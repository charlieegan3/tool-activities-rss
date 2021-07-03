package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charlieegan3/activities-rss/pkg/activity"
	"github.com/charlieegan3/activities-rss/pkg/download"
	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	email := viper.GetString("email")
	password := viper.GetString("password")
	baseURL := viper.GetString("baseURL")
	host := viper.GetString("host")
	feedTitle := viper.GetString("title")

	activities, err := download.Download(host, email, password)
	if err != nil {
		log.Fatalf("failed to download activities: %s", err)
	}

	feed := &feeds.Feed{
		Title:   feedTitle,
		Link:    &feeds.Link{Href: baseURL},
		Created: time.Now().UTC(),
	}

	groups := make(map[string][]activity.Activity)

	for _, activity := range activities {
		if activity.IsOnDate(time.Now().UTC()) {
			continue // will report tomorrow
		}

		groupKey, valid := activity.DateString()
		if !valid {
			continue
		}

		groups[groupKey] = append(groups[groupKey], activity)
	}

	for date, group := range groups {
		html := ""
		for _, activity := range group {
			html += activity.ToHTML(host)
		}

		item := feeds.Item{
			Title:       date,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/items/%s", baseURL, date)},
			Description: html,
			Id:          fmt.Sprintf("%s/items/%s", baseURL, date),
		}
		feed.Items = append(feed.Items, &item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("feed.xml", []byte(atom), 0644)
}
