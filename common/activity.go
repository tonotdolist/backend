package common

import (
	"time"
	"tonotdolist/internal/model"
)

type Activity struct {
	ActivityId  string
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time

	Completed bool `gorm:"default:false"`
}

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

type ActivityFetchByCountResponse struct {
	Activities []*Activity
}

type ActivityFetchByTimeRangeResponse struct {
	Activities []*Activity
}

func ConvertActivityDbModel(activity *model.Activity) *Activity {
	return &Activity{
		ActivityId:  activity.ActivityId,
		Type:        activity.Type,
		Name:        activity.Name,
		Priority:    activity.Priority,
		Description: activity.Description,
		Location:    activity.Location,
		Date:        activity.Date,
		Completed:   activity.Completed,
	}
}
