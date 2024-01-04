package models

type Schedule struct {
	Days []DaySchedule
}

type DaySchedule struct {
	DayOfTheWeek string
	Time         int64
	PoolTitle    string
	PoolVariant  []string
}

type ScheduleTime struct {
	GroupId int64
	Time    int64
}
