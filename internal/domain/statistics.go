package domain

type Statistics struct {
	TotalUsers       uint
	RegisteredUsers  uint
	OptedOutUsers    uint
	RegisteredDaily  uint
	RegisteredWeekly uint
	MaleCount        uint
	FemaleCount      uint
	Meetings         MeetingStats
}
