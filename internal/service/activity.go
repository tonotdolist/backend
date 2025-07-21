package service

import "tonotdolist/internal/repository"

type ActivityService interface {
}

type activityService struct {
	activityRepository repository.ActivityRepository
}

func NewActivityService(activityRepository repository.ActivityRepository) ActivityService {
	return &activityService{activityRepository: activityRepository}
}
