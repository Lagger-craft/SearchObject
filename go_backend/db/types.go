package db

import "searchobject/models"

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func parseTime(s string) models.Time {
	var t models.Time
	t.Scan(s)
	return t
}
