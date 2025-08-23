package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tonotdolist/common"
	"tonotdolist/internal/model"
	"tonotdolist/internal/repository"
)

type ActivityService interface {
	Create(ctx context.Context, userId string, request *common.ActivityCreateRequest) error
	Update(ctx context.Context, userId string, request *common.ActivityUpdateRequest) error
	Delete(ctx context.Context, userId string, request *common.ActivityDeleteRequest) error
	FetchByCount(ctx context.Context, userId string, request *common.ActivityFetchByCountRequest) (error, *common.ActivityFetchByCountResponse)
	FetchByTimeRange(ctx context.Context, userId string, request *common.ActivityFetchByTimeRangeRequest) (error, *common.ActivityFetchByTimeRangeResponse)
}
type activityService struct {
	activityRepository repository.ActivityRepository
}

func NewActivityService(activityRepository repository.ActivityRepository) ActivityService {
	return &activityService{activityRepository: activityRepository}
}

func (a *activityService) Create(ctx context.Context, userId string, request *common.ActivityCreateRequest) error {
	activityId, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating activity uuid: %w", err)
	}

	activityModel := &model.Activity{
		ActivityId:  activityId.String(),
		UserId:      userId,
		Type:        request.Type,
		Name:        request.Name,
		Priority:    request.Priority,
		Description: request.Description,
		Location:    request.Location,
		Date:        request.Date,
	}

	err = a.activityRepository.CreateActivity(ctx, activityModel)
	if err != nil {
		return fmt.Errorf("error creating activity in db:  %w", err)
	}

	return nil
}

func (a *activityService) Update(ctx context.Context, userId string, request *common.ActivityUpdateRequest) error {
	activityModel := &model.Activity{
		UserId:      userId,
		Type:        request.Type,
		Name:        request.Name,
		Priority:    request.Priority,
		Description: request.Description,
		Location:    request.Location,
		Date:        request.Date,
	}

	err := a.activityRepository.UpdateActivity(ctx, activityModel)
	if err != nil {
		return fmt.Errorf("error updating activity in db:  %w", err)
	}

	return nil
}

func (a *activityService) Delete(ctx context.Context, userId string, request *common.ActivityDeleteRequest) error {
	err := a.activityRepository.DeleteActivity(ctx, request.ActivityId, userId)
	if err != nil {
		return fmt.Errorf("error deleting activity in db:  %w", err)
	}

	return nil
}

func (a *activityService) FetchByCount(ctx context.Context, userId string, request *common.ActivityFetchByCountRequest) (error, *common.ActivityFetchByCountResponse) {
	activitiesModels, err := a.activityRepository.GetNUserActivity(ctx, userId, request.Count, request.Offset)
	if err != nil {
		return fmt.Errorf("error fetching activities by time range from db: %w", err), nil
	}

	return nil, &common.ActivityFetchByCountResponse{
		Activities: convertListToCommon(activitiesModels),
	}
}

func (a *activityService) FetchByTimeRange(ctx context.Context, userId string, request *common.ActivityFetchByTimeRangeRequest) (error, *common.ActivityFetchByTimeRangeResponse) {
	activitiesModels, err := a.activityRepository.GetUserActivityInRange(ctx, userId, request.Start, request.End)
	if err != nil {
		return fmt.Errorf("error fetching activities by time range from db: %w", err), nil
	}

	return nil, &common.ActivityFetchByTimeRangeResponse{
		Activities: convertListToCommon(activitiesModels),
	}
}

func convertListToCommon(activities []*model.Activity) []*common.Activity {
	internalActivities := make([]*common.Activity, len(activities))

	for i, activity := range activities {
		internalActivities[i] = common.ConvertActivityDbModel(activity)
	}

	return internalActivities
}
