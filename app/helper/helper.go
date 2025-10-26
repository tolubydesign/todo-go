package helper

import "time"

const (
	SuccessfulResponseMessage = "request handled successfully"
)

func IsStringUTC(date string, format string) (time.Time, error) {
	t, err := time.Parse(format, date)
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
