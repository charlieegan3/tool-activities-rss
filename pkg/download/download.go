package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/charlieegan3/activities-rss/pkg/activity"
	"github.com/charlieegan3/activities-rss/pkg/parser"
	"github.com/charlieegan3/activities-rss/pkg/types"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Download(host string, email, password string, athleteID string) ([]activity.Activity, error) {
	activities := []activity.Activity{}

	sessionURL := fmt.Sprintf("https://%s/session", host)
	loginURL := fmt.Sprintf("https://%s/login", host)

	// get login form token
	client := &http.Client{}
	req, err := http.NewRequest("GET", loginURL, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("DNT", "1")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.7,en-US;q=0.3")

	res, err := client.Do(req)
	if err != nil {
		return activities, err
	}

	cookie := fmt.Sprintf("%s", res.Header.Get("Set-Cookie"))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return activities, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	var authenticityToken string
	doc.Find(`form input[name=authenticity_token]`).Each(func(i int, s *goquery.Selection) {
		var found bool
		authenticityToken, found = s.Attr("value")
		if !found {
			log.Fatal("authenticity token not found")
		}
	})

	// get session token
	data := url.Values{}
	data.Add("utf8", "✓")
	data.Add("authenticity_token", authenticityToken)
	data.Add("plan", "")
	data.Add("email", email)
	data.Add("password", password)
	data.Add("remember_me", "on")

	req, err = http.NewRequest("POST", sessionURL, strings.NewReader(data.Encode()))

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("DNT", "1")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.7,en-US;q=0.3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err = client.Do(req)
	if err != nil {
		return activities, err
	}

	cookie = fmt.Sprintf("%s", res.Header.Get("Set-Cookie"))

	if res.StatusCode != http.StatusFound {
		log.Fatal("login attempt was not correctly redirected")
	}
	location := res.Header.Get("Location")
	if location != fmt.Sprintf("https://%s/dashboard", host) {
		log.Fatal("login attempt was not correctly redirect to dashboard")
	}

	// get dash
	var allEntries []types.RawEntry
	var before, cursor string
	before = "1656440258"
	cursor = "1656440448.929174"
	client = &http.Client{}
	for i := 0; i < 4; i++ {
		url := fmt.Sprintf(
			"https://%s/dashboard/feed?feed_type=following&athlete_id=%s&before=%s&cursor=%s",
			host, athleteID, before, cursor,
		)
		if i == 0 {
			url = fmt.Sprintf(
				"https://%s/dashboard/feed?feed_type=following&athlete_id=%s",
				host, athleteID,
			)
		}

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		req.Header.Add("Cookie", cookie)

		res, err = client.Do(req)
		if err != nil {
			return activities, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return activities, err
		}

		wrapper := struct {
			Entries []types.RawEntry `json:"entries"`
		}{}
		err = json.Unmarshal(body, &wrapper)
		if err != nil {
			return activities, err
		}

		allEntries = append(allEntries, wrapper.Entries...)

		last := wrapper.Entries[len(wrapper.Entries)-1]
		cursor = fmt.Sprintf("%f", last.CursorData.Rank)
		before = fmt.Sprintf("%d", last.CursorData.UpdatedAt)
	}

	activities, err = parser.Load(time.Now(), allEntries)
	if err != nil {
		return activities, fmt.Errorf("failed to parse collected raw activities: %w", err)
	}
	return activities, nil
}
