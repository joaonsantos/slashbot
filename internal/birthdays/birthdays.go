package birthdays

import (
	"errors"
	"strconv"
	"time"
)

var months = map[string]time.Month{
	"January":   time.January,
	"February":  time.February,
	"March":     time.March,
	"April":     time.April,
	"May":       time.May,
	"June":      time.June,
	"July":      time.July,
	"August":    time.August,
	"September": time.September,
	"October":   time.October,
	"November":  time.November,
	"December":  time.December,
}

type BirthdayEntry struct {
	userID string     `json: "user_id"`
	day    int        `json: "day"`
	month  time.Month `json: "month"`
}

func NewBirthdayEntry(userID string, day string, month string) (BirthdayEntry, error) {
	monthNumber, ok := months[month]
	if !ok {
		return BirthdayEntry{}, errors.New("please input the long form of the month, eg. \"January\"")
	}
	dayNumber, err := strconv.Atoi(day)
	if err != nil {
		return BirthdayEntry{}, errors.New("please input the day as a number, eg. \"1\"")
	}
	return BirthdayEntry{
		userID: userID,
		day:    dayNumber,
		month:  monthNumber,
	}, nil
}
