package tool

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"
)

// ActivitiesRSS is a tool for creating RSS feeds from run/bike/swim etc activities found on a certain website
type ActivitiesRSS struct {
	config *gabs.Container
	db     *sql.DB
}

func (d *ActivitiesRSS) Name() string {
	return "activities-rss"
}

func (d *ActivitiesRSS) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config: true,
		Jobs:   true,
	}
}

func (d *ActivitiesRSS) SetConfig(config map[string]any) error {
	d.config = gabs.Wrap(config)

	return nil
}
func (d *ActivitiesRSS) Jobs() ([]apis.Job, error) {
	var j []apis.Job
	var path string
	var ok bool

	// load all config
	path = "jobs.new-entry.schedule"
	schedule, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.new-entry.endpoint"
	endpoint, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.new-entry.host"
	host, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.new-entry.athlete_id"
	athleteID, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.new-entry.email"
	email, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.new-entry.password"
	password, ok := d.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	return []apis.Job{
		&NewEntry{
			ScheduleOverride: schedule,
			Endpoint:         endpoint,

			Host:      host,
			AthleteID: athleteID,
			Email:     email,
			Password:  password,
		},
	}, nil
}

func (d *ActivitiesRSS) DatabaseMigrations() (*embed.FS, string, error) {
	return &embed.FS{}, "migrations", nil
}
func (d *ActivitiesRSS) DatabaseSet(db *sql.DB)              {}
func (d *ActivitiesRSS) HTTPPath() string                    { return "" }
func (d *ActivitiesRSS) HTTPAttach(router *mux.Router) error { return nil }
