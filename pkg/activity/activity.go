package activity

import (
	"fmt"
	"strings"
	"time"
)

type Activity struct {
	Name  string
	Date  string
	ID    string
	Title string
	Type  string
	Stats []string
}

func (a *Activity) ToHTML(host string) string {
	statsHTML := ""
	for _, stat := range a.Stats {
		parts := strings.Split(stat, "\n")
		if len(parts) > 1 {
			statsHTML += fmt.Sprintf(`<li><b>%s</b> %s</li>`, parts[0], parts[1])
		}
	}
	tmpl := `<div style="font-size: 0.8rem; font-family: sans-serif;">
  <h3><a href="https://%s/activities/%s">%s (%s)</a></h3>
  <h4> %s </h4>
  <ul> %s </ul>
</div>
<hr/>`

	return fmt.Sprintf(tmpl, host, a.ID, a.Name, a.Type, a.Title, statsHTML)
}

func (a *Activity) IsOnDate(date time.Time) bool {
	return strings.HasPrefix(a.Date, date.Format("2006-01-02"))
}

func (a *Activity) DateString() (string, bool) {
	parts := strings.Split(a.Date, " ")
	if len(parts) > 0 {
		return parts[0], true
	}
	return "", false
}
