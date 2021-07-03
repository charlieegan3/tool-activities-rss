package download

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charlieegan3/activities-rss/pkg/activity"
	"github.com/gocolly/colly"
)

func Download(host string, email, password string) ([]activity.Activity, error) {
	activities := []activity.Activity{}

	sessionURL := fmt.Sprintf("https://%s/session", host)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Accept-Language", "en-GB,en;q=0.7,en-US;q=0.3")
	})

	var authenticityToken string
	c.OnHTML(`input[name="authenticity_token"]`, func(e *colly.HTMLElement) {
		authenticityToken = e.Attr("value")
	})

	c.Visit(fmt.Sprintf("https://%s/login", host))

	if authenticityToken == "" {
		return activities, fmt.Errorf("failed to get form authenticity token")
	}

	c.OnHTML(`.feed-entry`, func(e *colly.HTMLElement) {
		activity := activity.Activity{
			Name:  e.DOM.Find(".entry-owner").Text(),
			Title: strings.TrimSpace(e.DOM.Find(".title-text").Text()),
		}
		path, _ := e.DOM.Find(".title-text a").Attr("href")
		pathComponents := strings.Split(path, "/")
		if len(pathComponents) > 0 {
			activity.ID = pathComponents[len(pathComponents)-1]
		}

		date, exists := e.DOM.Find(".timestamp").Attr("datetime")
		if exists {
			activity.Date = date
		}

		html, err := e.DOM.Find(".entry-icon").Html()
		if err == nil {
			re := regexp.MustCompile(` icon-\w+`)
			icons := re.FindAllString(html, 1)
			if len(icons) == 1 {
				activity.Type = strings.TrimPrefix(icons[0], " icon-")
			}
		}

		e.DOM.Find(".list-stats .stat").Each(func(i int, s *goquery.Selection) {
			activity.Stats = append(activity.Stats, strings.TrimSpace(s.Text()))
		})

		activities = append(activities, activity)
	})

	data := url.Values{}
	data.Add("utf8", "✓")
	data.Add("authenticity_token", authenticityToken)
	data.Add("plan", "")
	data.Add("email", email)
	data.Add("password", password)
	data.Add("remember_me", "on")

	c.Request("POST", sessionURL, strings.NewReader(data.Encode()), &colly.Context{}, http.Header{})

	return activities, nil
}
