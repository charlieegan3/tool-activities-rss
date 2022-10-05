package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charlieegan3/tool-activities-rss/pkg/activity"
	"github.com/charlieegan3/tool-activities-rss/pkg/download"
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
	athleteID := viper.GetString("athlete_id")
	baseURL := viper.GetString("baseURL")
	host := viper.GetString("host")
	feedTitle := viper.GetString("title")

	activities, err := download.Download(host, email, password, athleteID)
	if err != nil {
		log.Fatalf("failed to download activities: %s", err)
	}

	fmt.Println(len(activities))

	feed := &feeds.Feed{
		Title:   feedTitle,
		Link:    &feeds.Link{Href: baseURL},
		Created: time.Now().UTC(),
	}

	groups := make(map[string][]activity.Activity)

	for _, activity := range activities {
		if !activity.IsOnDate(time.Now().Add(-24 * time.Hour)) {
			continue // only report yesterday
		}

		groupKey := activity.DateString()

		groups[groupKey] = append(groups[groupKey], activity)
	}

	for date, group := range groups {
		html := ""
		for _, activity := range group {
			html += activity.ToHTML(host)
		}

		item := feeds.Item{
			Title:       date,
			Updated:     time.Now(),
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
