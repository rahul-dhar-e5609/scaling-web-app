package main

import "github.com/IAmRDhar/scaling-web-app/backend/entity"

type logEntries []entity.LogEntry

func (le logEntries) Len() int {
	return len(le)
}

func (le logEntries) Swap(i, j int) {
	le[i], le[j] = le[j], le[i]
}

func (le logEntries) Less(i, j int) bool {
	return le[i].Timestamp.Before(le[j].Timestamp)
}
