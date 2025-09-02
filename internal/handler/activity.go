package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"tonotdolist/api"
	"tonotdolist/common"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/log"
)

var (
	activityCreateRequestType = reflect.TypeOf(common.ActivityCreateRequest{})
)

type ActivityHandler struct {
	*Handler
	activityService service.ActivityService
}

func (h *ActivityHandler) CreateActivity(ctx *gin.Context) {
	userId := getUserId(ctx)

	logger := log.GetLoggerFromContext(ctx)
	rawReq, err := api.BindJSON(ctx, activityCreateRequestType)

	defer api.HandleResponse(ctx, err, nil)

	if err != nil {
		if !errors.Is(err, common.ErrBadRequest) {
			logger.Error().Err(err).Msg("unable to bind json")
		}

		return
	}

	req := rawReq.(*common.ActivityCreateRequest)

	err = h.activityService.Create(ctx, userId, req)
	if err != nil {
		logger.Error().Err(err).Msg("error creating activity")
	}
}
