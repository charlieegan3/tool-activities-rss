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
				Name:        "Morning Run",
				AthleteName: "Joe Bloggs",
				Description: "desc",
				Date:        time.Date(2022, 6, 28, 10, 11, 12, 0, time.UTC),
				ID:          "1",
				Type:        "run",
				Location:    "Location",
				Stats:       []string{"Time\n1h 31m"},
			},
			ExpectedHTML: `<div style="font-size: 0.8rem; font-family: sans-serif;">
  <h3><a href="https://example.com/activities/1">Joe Bloggs (run)</a></h3>
  <h4> Morning Run (10:11, Location)</h4>
  <div> desc </div>
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
