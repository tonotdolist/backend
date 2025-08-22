package common

import "time"

type ActivityCreateRequest struct {
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time
}

type ActivityUpdateRequest struct {
	ActivityId  string
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time

	Completed bool `gorm:"default:false"`
}

type ActivityDeleteRequest struct {
	ActivityId string
}

type ActivityFetchByCountRequest struct {
	Count int
}

type ActivityFetchByTimeRangeRequest struct {
	Start time.Time
	End   time.Time
}
