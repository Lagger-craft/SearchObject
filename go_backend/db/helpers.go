package db

import (
	"database/sql"
	"strconv"
)

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nullFromSQL(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func isUnique(err error) bool {
	return err != nil && (contains(err.Error(), "UNIQUE constraint") || contains(err.Error(), "UNIQUE"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
