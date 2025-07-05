package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	iapi "tonotdolist/api"
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
	var resp *common.UserLoginResponse

	logger := log.GetLoggerFromContext(ctx)
	rawReq, err := api.BindJSON(ctx, loginRequestType)

	defer func() {
		iapi.HandleResponse(ctx, err, resp)
	}()

	if err != nil {
		if !errors.Is(err, common.ErrBadRequest) {
			logger.Error().Err(err).Msg("unable to bind json")
		}

		return
	}

	req := rawReq.(*common.UserLoginRequest)
	sessionId, err := h.userService.Login(ctx, req)

	if err != nil {
		if !common.IsCommonError(err) {
			logger.Error().Err(err).Msg("error handling login request")
		}

		return
	}

	resp = &common.UserLoginResponse{
		SessionID: sessionId,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var resp *common.UserRegisterResponse

	logger := log.GetLoggerFromContext(ctx)

	rawReq, err := api.BindJSON(ctx, registerRequestType)

	defer func() {
		iapi.HandleResponse(ctx, err, resp)
	}()

	if err != nil {
		if !errors.Is(err, common.ErrBadRequest) {
			logger.Error().Err(err).Msg("unable to bind json")
		}

		return
	}

	req := rawReq.(*common.UserRegisterRequest)
	sessionId, err := h.userService.Register(ctx, req)

	if err != nil {
		if !common.IsCommonError(err) {
			logger.Error().Err(err).Msg("error handling register request")
		}
	}

	resp = &common.UserRegisterResponse{
		SessionID: sessionId,
	}
}
