package helper

import "time"

const (
	SuccessfulResponseMessage = "request handled successfully"
)

/*
Find out if the provided string can be converted into a "time".

Example:

date = 2006-01-02T15:04:05Z07:00 OR 2026-01-30T19:00:00Z

layout = time.RFC3339
*/
func IsStringUTC(date string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, date)
	if err != nil {
		return t, err
	}
	// return t.Location() == time.UTC, nil
	return t, nil
}

func ConvertIntToPointerInt(value int) *int {
	return &value
}

func ConvertTimeToString(t time.Time) string {
	dueDate := ""
	dueDate = t.Format(time.RFC3339)
	return dueDate
}
