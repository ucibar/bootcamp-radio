package handler

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/uCibar/bootcamp-radio/api/response"
	"github.com/uCibar/bootcamp-radio/entity"
	"github.com/uCibar/bootcamp-radio/service/auth"
	"log"
	"net/http"
)

type AuthHandler struct {
	*handler
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service, logger *log.Logger) *AuthHandler {
	return &AuthHandler{handler: newHandler(logger), authService: authService}
}

func (h *AuthHandler) Middleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			h.writeResponse(w, response.InvalidCredentials("Authorization header is required"))
			return
		}

		u, err := h.authService.VerifyToken(authToken)
		if err != nil {
			h.writeLog(err.Error())
			h.writeResponse(w, response.InvalidCredentials("token is invalid"))
			return
		}

		c := context.WithValue(r.Context(), "user", u)

		handler(w, r.WithContext(c), ps)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var authReq entity.Auth

	err := h.readBody(r, &authReq)
	if err != nil {
		h.writeResponse(w, response.InvalidBody(err.Error()))
		return
	}

	token, err := h.authService.AuthenticateWithToken(&authReq)
	if errors.Is(err, entity.ErrUserNotFound) || errors.Is(err, entity.ErrPasswordIncorrect) {
		h.writeResponse(w, response.InvalidCredentials("invalid username or password"))
		return
	} else if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(http.StatusOK, token))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var registerReg entity.Register

	err := h.readBody(r, &registerReg)
	if err != nil {
		h.writeResponse(w, response.InvalidBody(err.Error()))
		return
	}

	err = h.authService.RegisterUser(&registerReg)
	if errors.Is(err, entity.ErrPasswordRepeatIncorrect) || errors.Is(err, entity.ErrEmailExist) || errors.Is(err, entity.ErrUsernameExist) {
		h.writeResponse(w, response.BadRequest(err.Error()))
		return
	} else if err != nil {
		h.writeLog(err.Error())
		h.writeResponse(w, response.ServerError())
		return
	}

	h.writeResponse(w, response.Success(http.StatusOK, nil))
}
