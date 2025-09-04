package utils

import (
	"database/sql"
	"time"
)

// ToNullString converts a *string into sql.NullString.
//
// If the pointer is nil, it returns a sql.NullString with Valid = false,
// which will be treated as NULL in the database.
// If the pointer is not nil, it returns a sql.NullString with the string
// value set and Valid = true.
func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// ToNullInt32 converts a *int into sql.NullInt32.
//
// If the pointer is nil, it returns a sql.NullInt32 with Valid = false,
// which will be treated as NULL in the database.
// If the pointer is not nil, it returns a sql.NullInt32 with the int32
// value set and Valid = true.
func ToNullInt32(i *int) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

// ToNullTime converts a *time.Time into sql.NullTime.
//
// If the pointer is nil, it returns a sql.NullTime with Valid = false,
// which will be treated as NULL in the database.
// If the pointer is not nil, it returns a sql.NullTime with the time
// value set and Valid = true.
func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// ToNullTimeFromString converts a *string into sql.NullTime, assuming
// the string is in RFC3339 format (ISO8601).
//
// If the pointer is nil, it returns a sql.NullTime with Valid = false.
// If the string is invalid (fails to parse), it returns an error.
// If valid, it returns a sql.NullTime with the parsed time and Valid = true.
func ToNullTimeFromString(s *string) (sql.NullTime, error) {
	if s == nil {
		return sql.NullTime{Valid: false}, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

// NullTimeToString converts a sql.NullTime into a string.
//
// If the NullTime is valid, it returns the time formatted as RFC3339.
// If not valid (NULL in DB), it returns an empty string.
func NullTimeToString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}
