package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/jwt"
	"github.com/sainak/bitsb/utils/response"
)

type UserHandler struct {
	service domain.UserServiceProvider
}

func New(service domain.UserServiceProvider) *UserHandler {
	return &UserHandler{service}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := &domain.UserLogin{}
	err := render.Bind(r, data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("empty request body")
		}
		response.RespondForError(w, r, err)
		return
	}

	token, err := u.service.Login(r.Context(), data)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, token)
}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	data := &domain.RefreshToken{}
	err := render.Bind(r, data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("empty request body")
		}
		response.RespondForError(w, r, err)
		return
	}

	token, err := u.service.RefreshToken(r.Context(), data.RefreshToken)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, token)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	data := &domain.UserRegister{}
	err := render.Bind(r, data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("empty request body")
		}
		response.RespondForError(w, r, err)
		return
	}

	user := &domain.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	err = u.service.Signup(r.Context(), user)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, user)
}

func (u *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserID).(int64)
	user, err := u.service.GetUserByID(r.Context(), userID)
	if err != nil {
		response.RespondForError(w, r, err)
		return
	}
	render.JSON(w, r, user)
}
