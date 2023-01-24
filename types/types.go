package types

import "time"

type Player struct {
	Position  int64
	FirstName string
}

type PastMatch struct {
	Player1        string
	Player2        string
	Winner         int
	Player1PrevPos int
	Player2PrevPos int
	Date           time.Time
}