package parser

import (
	"fmt"
	"github.com/charlieegan3/tool-activities-rss/pkg/activity"
	"github.com/charlieegan3/tool-activities-rss/pkg/types"
	"strings"
	"time"
)

func Load(collectTime time.Time, entries []types.RawEntry) ([]activity.Activity, error) {
	var activities []activity.Activity

	for _, e := range entries {
		if e.Entity != "Activity" {
			continue
		}

		activity := activity.Activity{
			ID: e.Activity.ID,

			Type: e.Activity.Type,

			Name:        e.Activity.ActivityName,
			Description: e.Activity.Description,
			Location:    e.Activity.TimeAndLocation.Location,

			AthleteName: e.Activity.Athlete.AthleteName,

			MapURL: e.Activity.MapAndPhotos.ActivityMap.URL,
		}

		for _, p := range e.Activity.MapAndPhotos.PhotoList {
			if p.Large != "" {
				activity.PhotoURL = p.Large
				break
			}
		}

		var stats []string
		statValues := make(map[string]string)
		statNames := make(map[string]string)
		for _, s := range e.Activity.Stats {
			key := strings.TrimSuffix(s.Key, "_subtitle")

			if strings.HasSuffix(s.Key, "_subtitle") {
				stats = append(stats, key)
				statNames[key] = s.Value
			} else {
				statValues[key] = s.Value
			}
		}

		for _, v := range stats {
			name, ok := statNames[v]
			if !ok {
				continue
			}
			value, ok := statValues[v]
			if !ok {
				continue
			}

			activity.Stats = append(activity.Stats, fmt.Sprintf("<b>%s</b>: %s", name, value))
		}

		activity.Date = collectTime
		if e.Activity.TimeAndLocation.DisplayDate == "Yesterday" {
			activity.Date = collectTime.Add(-24 * time.Hour)
			parsedTime, err := time.Parse("Yesterday at 3:04 PM", e.Activity.TimeAndLocation.DisplayDateAtTime)
			if err == nil {
				activity.Date, _ = time.Parse(
					"2006-01-02T15:04",
					fmt.Sprintf("%sT%s",
						activity.Date.Format("2006-01-02"),
						parsedTime.Format("03:04"),
					),
				)
			}
		} else if e.Activity.TimeAndLocation.DisplayDate == "Today" {
			parsedTime, err := time.Parse("Today at 3:04 PM", e.Activity.TimeAndLocation.DisplayDateAtTime)
			if err == nil {
				activity.Date, _ = time.Parse(
					"2006-01-02T15:04",
					fmt.Sprintf("%sT%s",
						activity.Date.Format("2006-01-02"),
						parsedTime.Format("03:04"),
					),
				)
			}
		} else {
			activity.Date, _ = time.Parse("2 January 2006 at 15:04", e.Activity.TimeAndLocation.DisplayDateAtTime)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}
