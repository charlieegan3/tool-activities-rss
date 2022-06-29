package parser

import (
	"encoding/json"
	"github.com/charlieegan3/activities-rss/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// fixture file generated with:
	// cat raw.json | gron | grep 'activityName\|athleteName\|stats\|activity\.id\|photoList\[\d*\]\.large\|activityMap.url\|timeAndLocation.displayDateAtTime\|activity.type' | gron --ungron > fixture2.json
	data, err := ioutil.ReadFile("fixture.json")
	require.NoError(t, err)

	wrapper := struct {
		Entries []types.RawEntry `json:"entries"`
	}{}
	err = json.Unmarshal(data, &wrapper)
	require.NoError(t, err)

	collectTime := time.Date(2022, 6, 29, 12, 25, 0, 0, time.UTC)

	activities, err := Load(collectTime, wrapper.Entries)
	require.NoError(t, err)

	assert.Equal(t, 19, len(activities))
	assert.Equal(
		t,
		"Wahoo SYSTM: Provence: Gorges de la Nesque",
		activities[0].Name,
	)
	assert.Equal(t, "NAME", activities[0].AthleteName)
	assert.Equal(t, "Ride", activities[0].Type)
	assert.Equal(t, "<p>desc</p>", activities[0].Description)
	assert.Equal(
		t,
		"https://dgtzuqphqg23d.cloudfront.net/R8w_cwgwN5Cophyx_j2aRMoSw1Ak1pdK1f3VS7Pz1no-2048x1365.jpg",
		activities[0].PhotoURL,
	)
	assert.Equal(
		t,
		"https://d3o5xota0a1fcr.cloudfront.net/v6/maps/FAMYRMVGLVLUWOJCVUBINWYSU4HGDMCJQLRF3VBAG32GSEQGVA5757AJVFHXDLWNA5ZUUYAIUF3DT3AYDEALO2O2UQXJ4HF6",
		activities[1].MapURL,
	)
	assert.Equal(t, "South Lake Tahoe, USA", activities[0].Location)
	assert.Equal(t,
		time.Date(2022, 6, 28, 12, 8, 0, 0, time.UTC),
		activities[0].Date,
	)

	assert.Equal(
		t,
		[]string{
			"<b>Distance</b>: 27.00<abbr class='unit' title='kilometers'> km</abbr>",
			"<b>Time</b>: 50<abbr class='unit' title='minute'>m</abbr> 57<abbr class='unit' title='second'>s</abbr>",
			"<b>Avg Power</b>: 206<abbr class='unit' title='watts'> W</abbr>",
		},
		activities[0].Stats,
	)
	assert.Equal(t,
		time.Date(2022, 6, 27, 18, 7, 0, 0, time.UTC),
		activities[1].Date,
	)
}
