package v1

import (
	"time"
	"tonotdolist/common"
	"tonotdolist/pkg/api"
)

func init() {
	api.RegisterRequest[common.UserRegisterRequest, UserRegisterRequest](version)
	api.RegisterRequest[common.UserLoginRequest, UserLoginRequest](version)
}

type ActivityCreateRequest struct {
	ActivityId  string
	UserId      string `gorm:"not null"`
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time
}

func (r *ActivityCreateRequest) ToInternalRequest() interface{} {
	return &common.ActivityCreateRequest{
		ActivityId:  r.ActivityId,
		UserId:      r.UserId,
		Type:        r.Type,
		Name:        r.Name,
		Priority:    r.Priority,
		Description: r.Description,
		Location:    r.Location,
		Date:        r.Date,
	}
}

type ActivityUpdateRequest struct {
	ActivityId  string
	UserId      string `gorm:"not null"`
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time

	Completed bool `gorm:"default:false"`
}

func (r *ActivityUpdateRequest) ToInternalRequest() interface{} {
	return &common.ActivityUpdateRequest{
		ActivityId:  r.ActivityId,
		UserId:      r.UserId,
		Type:        r.Type,
		Name:        r.Name,
		Priority:    r.Priority,
		Description: r.Description,
		Location:    r.Location,
		Date:        r.Date,

		Completed: r.Completed,
	}
}

type ActivityDeleteRequest struct {
	ActivityId string
}

func (r *ActivityDeleteRequest) ToInternalRequest() interface{} {
	return &common.ActivityDeleteRequest{
		ActivityId: r.ActivityId,
	}
}

type ActivityFetchByCountRequest struct {
	Count int
}

func (r *ActivityFetchByCountRequest) ToInternalRequest() interface{} {
	return &common.ActivityFetchByCountRequest{
		Count: r.Count,
	}
}

type ActivityFetchByTimeRangeRequest struct {
	Start time.Time
	End   time.Time
}

func (r *ActivityFetchByTimeRangeRequest) ToInternalRequest() interface{} {
	return &common.ActivityFetchByTimeRangeRequest{
		Start: r.Start,
		End:   r.End,
	}
}

var _ api.VersionedRequest = (*ActivityCreateRequest)(nil)
var _ api.VersionedRequest = (*ActivityUpdateRequest)(nil)
var _ api.VersionedRequest = (*ActivityDeleteRequest)(nil)
var _ api.VersionedRequest = (*ActivityFetchByCountRequest)(nil)
var _ api.VersionedRequest = (*ActivityFetchByTimeRangeRequest)(nil)
