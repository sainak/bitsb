package handler

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/api"
	"github.com/sainak/bitsb/domain/middleware"
)

type UserHandler struct {
	service domain.UserServiceProvider
}

func New(service domain.UserServiceProvider) *UserHandler {
	return &UserHandler{service}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := &domain.UserLoginForm{}
	err := render.Bind(r, data)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	token, err := u.service.Login(r.Context(), data)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, token)
}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	data := &domain.RefreshTokenFrom{}
	err := render.Bind(r, data)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	token, err := u.service.RefreshToken(data.RefreshToken)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, token)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	data := &domain.UserRegisterForm{}
	err := render.Bind(r, data)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}

	user := &domain.User{
		FirstName:      data.FirstName,
		LastName:       data.LastName,
		Email:          data.Email,
		Password:       data.Password,
		Access:         domain.Passenger,
		HomeLocationID: data.HomeLocationID,
		WorkLocationID: data.WorkLocationID,
	}

	err = u.service.Signup(r.Context(), user)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (u *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// extract the user from the context
	user := r.Context().Value(middleware.UserCtxKey).(*domain.User)
	render.JSON(w, r, user)
}
