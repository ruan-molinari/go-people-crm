package model

import "time"

type Person struct {
	ID int64
	Name string
	Phone string
	LastTiemMet time.Time
	MeetingFrequecy time.Duration
}

func (p Person) ShouldBeMetAgainIn() time.Time {
	return p.LastTiemMet.Add(p.MeetingFrequecy)
}
