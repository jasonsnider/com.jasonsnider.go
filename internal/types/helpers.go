package types

import (
	"database/sql"
	"fmt"
	"time"
)

// Helper function to convert string to sql.NullString
func TypeSqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Helper function to convert time.Time to sql.NullTime
func TypeSqlNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Time: t, Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// Helper function to convert string in SQL format to sql.NullTime
func ParseSqlNullTime(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{Valid: false}, nil // Treat empty string as NULL
	}

	// Define the date format matching SQL standard: "YYYY-MM-DD HH:mm:ss"
	const layout = "2006-01-02 15:04:05"

	// Parse the string into time.Time
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return sql.NullTime{}, err
	}

	// Return a valid sql.NullTime
	return sql.NullTime{
		Time:  parsedTime,
		Valid: true,
	}, nil
}

func SafeValue(value interface{}) string {
	switch v := value.(type) {
	case sql.NullString:
		if v.Valid {
			return v.String
		}
	case sql.NullInt64:
		if v.Valid {
			return fmt.Sprintf("%d", v.Int64)
		}
	case sql.NullFloat64:
		if v.Valid {
			return fmt.Sprintf("%f", v.Float64)
		}
	case sql.NullBool:
		if v.Valid {
			return fmt.Sprintf("%t", v.Bool)
		}
	case sql.NullTime:
		if v.Valid {
			return v.Time.Format("2006-01-02 15:04:05") // Custom time format
		}
	}
	return ""
}
