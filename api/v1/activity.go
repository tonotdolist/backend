package v1

import (
	"time"
	"tonotdolist/common"
	"tonotdolist/pkg/api"
)

func init() {
	api.RegisterRequest[common.ActivityCreateRequest, ActivityCreateRequest](version)
	api.RegisterRequest[common.ActivityDeleteRequest, ActivityDeleteRequest](version)
	api.RegisterRequest[common.ActivityUpdateRequest, ActivityUpdateRequest](version)
	api.RegisterRequest[common.ActivityFetchByCountRequest, ActivityFetchByCountRequest](version)
	api.RegisterRequest[common.ActivityFetchByTimeRangeRequest, ActivityFetchByTimeRangeRequest](version)
}

type ActivityCreateRequest struct {
	Type        int8      `json:"type" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Priority    int8      `json:"priority" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}

func (r *ActivityCreateRequest) ToInternalRequest() interface{} {
	return &common.ActivityCreateRequest{
		Type:        r.Type,
		Name:        r.Name,
		Priority:    r.Priority,
		Description: r.Description,
		Location:    r.Location,
		Date:        r.Date,
	}
}

type ActivityUpdateRequest struct {
	ActivityId  string    `json:"activity_id" binding:"required,id"`
	Type        int8      `json:"type" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Priority    int8      `json:"priority" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`

	Completed bool `json:"completed" binding:"required"`
}

func (r *ActivityUpdateRequest) ToInternalRequest() interface{} {
	return &common.ActivityUpdateRequest{
		ActivityId:  r.ActivityId,
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
	ActivityId string `json:"activity_id" binding:"required,id"`
}

func (r *ActivityDeleteRequest) ToInternalRequest() interface{} {
	return &common.ActivityDeleteRequest{
		ActivityId: r.ActivityId,
	}
}

type ActivityFetchByCountRequest struct {
	Count int `json:"activity_id" binding:"required"`
}

func (r *ActivityFetchByCountRequest) ToInternalRequest() interface{} {
	return &common.ActivityFetchByCountRequest{
		Count: r.Count,
	}
}

type ActivityFetchByTimeRangeRequest struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
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
