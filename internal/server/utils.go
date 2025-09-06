package server

import "database/sql"

func toString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func int32toInt(ni sql.NullInt32) int {
	if ni.Valid {
		return int(ni.Int32)
	}
	return 0
}
