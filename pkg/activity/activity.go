package activity

import (
	"fmt"
	"time"
)

type Activity struct {
	ID   string
	Type string

	AthleteName string

	Name        string
	Description string
	Location    string

	PhotoURL string
	MapURL   string

	Date time.Time

	Stats []string
}

func (a *Activity) ToHTML(host string) string {
	statsHTML := ""
	for _, stat := range a.Stats {
		statsHTML += fmt.Sprintf(`<li>%s</li>`, stat)
	}
	tmpl := `<div style="font-size: 0.8rem; font-family: sans-serif;">
  <h3><a href="https://%s/activities/%s">%s (%s)</a></h3>
  <h4> %s (%s, %s)</h4>
  <div> %s </div>
  <ul> %s </ul>
</div>
<hr/>`

	return fmt.Sprintf(
		tmpl,
		host,
		a.ID,
		a.AthleteName,
		a.Type,
		a.Name,
		a.Date.Format("15:04"),
		a.Location,
		a.Description,
		statsHTML,
	)
}

func (a *Activity) IsOnDate(date time.Time) bool {
	return a.Date.Format("2006-01-02") == date.Format("2006-01-02")
}

func (a *Activity) DateString() string {
	return a.Date.Format("2006-01-02")
}
