package types

type RawEntry struct {
	Entity   string `json:"entity"`
	Activity struct {
		ActivityName string `json:"activityName"`
		Athlete      struct {
			AthleteName string `json:"athleteName"`
		} `json:"athlete"`
		Description  string `json:"description"`
		ID           string `json:"id"`
		IsVirtual    bool   `json:"isVirtual"`
		MapAndPhotos struct {
			ActivityMap struct {
				URL string `json:"url"`
			} `json:"activityMap"`
			PhotoList []struct {
				Large string `json:"large"`
			} `json:"photoList"`
		} `json:"mapAndPhotos"`
		OwnedByCurrentAthlete bool `json:"ownedByCurrentAthlete"`
		Stats                 []struct {
			Key         string      `json:"key"`
			Value       string      `json:"value"`
			ValueObject interface{} `json:"value_object"`
		} `json:"stats"`
		TimeAndLocation struct {
			DisplayDate       string `json:"displayDate"`
			DisplayDateAtTime string `json:"displayDateAtTime"`
			Location          string `json:"location"`
			TimestampFormat   string `json:"timestampFormat"`
		} `json:"timeAndLocation"`
		Type        string      `json:"type"`
		Visibility  string      `json:"visibility"`
		WorkoutType interface{} `json:"workoutType"`
	} `json:"activity"`
	CursorData struct {
		Rank      float64 `json:"rank"`
		UpdatedAt int64   `json:"updated_at"`
	} `json:"cursorData"`
	TimeAndLocation struct {
		DisplayDate       string `json:"displayDate"`
		DisplayDateAtTime string `json:"displayDateAtTime"`
		Location          string `json:"location"`
		TimestampFormat   string `json:"timestampFormat"`
	} `json:"timeAndLocation"`
}
