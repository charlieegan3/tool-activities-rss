package activity

import (
	"testing"
	"time"
)

func TestActivityToHTML(t *testing.T) {
	testCases := []struct {
		Description  string
		Activity     Activity
		ExpectedHTML string
	}{
		{
			Description: "simple example",
			Activity: Activity{
				Name:  "Joe Bloggs",
				Date:  "2021-07-02 06:08:17 UTC",
				ID:    "1",
				Title: "Morning Run",
				Type:  "run",
				Stats: []string{"Time\n1h 31m"},
			},
			ExpectedHTML: `<div style="font-size: 0.8rem; font-family: sans-serif;">
  <h3><a href="https://example.com/activities/1">Joe Bloggs (run)</a></h3>
  <h4> Morning Run </h4>
  <ul> <li><b>Time</b> 1h 31m</li> </ul>
</div>
<hr/>`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			result := test.Activity.ToHTML("example.com")
			if result != test.ExpectedHTML {
				t.Fatalf("unexpected output: got:\n%s\n\nwant:\n%s", result, test.ExpectedHTML)
			}
		})
	}
}

func TestActivityIsOnDate(t *testing.T) {
	testCases := []struct {
		Description string
		Activity    Activity
		Date        time.Time
		Expected    bool
	}{
		{
			Description: "matches",
			Activity: Activity{
				Name: "Joe Bloggs",
				Date: "2021-07-02 06:08:17 UTC",
			},
			Expected: true,
			Date:     time.Date(2021, time.July, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Description: "does not match",
			Activity: Activity{
				Name: "Joe Bloggs",
				Date: "2021-07-03 06:08:17 UTC",
			},
			Expected: false,
			Date:     time.Date(2021, time.July, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			result := test.Activity.IsOnDate(test.Date)
			if result != test.Expected {
				t.Fatalf("unexpected output: got:\n%v\n\nwant:\n%v", result, test.Expected)
			}
		})
	}
}
