package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"tonotdolist/common"
	"tonotdolist/internal/service"
	"tonotdolist/pkg/api"
	"tonotdolist/pkg/log"
)

var (
	loginRequestType    = reflect.TypeOf(&common.UserLoginRequest{})
	registerRequestType = reflect.TypeOf(&common.UserRegisterRequest{})
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	logger := log.GetLoggerFromContext(ctx)
	version := ctx.GetUint("version")

	rawReq, err := api.BindJSON(ctx, version, loginRequestType)

	defer h.responder.HandleResponse(version, ctx, err, nil)

	if err != nil {
		if !errors.Is(err, common.ErrBadRequest) {
			logger.Error().Err(err).Msg("unable to bind json")
		}

		return
	}

	req := rawReq.(*common.UserLoginRequest)
	err = h.userService.Login(ctx, req)

	if err != nil {
		if !common.IsCommonError(err) {
			logger.Error().Err(err).Msg("error handling login request")
		}
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	logger := log.GetLoggerFromContext(ctx)
	version := ctx.GetUint("version")

	rawReq, err := api.BindJSON(ctx, version, registerRequestType)

	defer h.responder.HandleResponse(version, ctx, err, nil)

	if err != nil {
		if !errors.Is(err, common.ErrBadRequest) {
			logger.Error().Err(err).Msg("unable to bind json")
		}

		return
	}

	req := rawReq.(*common.UserRegisterRequest)
	err = h.userService.Register(ctx, req)

	if err != nil {
		if !errors.Is(err, common.ErrUnauthorized) {
			logger.Error().Err(err).Msg("error handling register request")
		}
	}
}
