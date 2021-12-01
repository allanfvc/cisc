package utils

import "time"

func WeekStartDate(date time.Time) time.Time {
  offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
  result := date.Add(time.Duration(offset*24) * time.Hour)
  return result
}
