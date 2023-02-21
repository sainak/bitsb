package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/sainak/bitsb/apperrors"
	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/mocks"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationHandlerTestSuite struct {
	suite.Suite
	handler *LocationHandler
	service *mocks.LocationServiceProvider
}

func TestLocationHandlerTestSuite(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func(patch *mpatch.Patch) {
		if err := patch.Unpatch(); err != nil {
			t.Fatal(err)
		}
	}(patch)

	suite.Run(t, new(LocationHandlerTestSuite))
}

func (s *LocationHandlerTestSuite) SetupTest() {
	s.service = new(mocks.LocationServiceProvider)
	s.handler = NewLocationHandler(s.service)
}

func (s *LocationHandlerTestSuite) TestListAll() {
	t := s.T()

	url := "/locations"

	locations := []*bitsb.Location{
		{ID: 1, Name: "Test Location 1"},
		{ID: 2, Name: "Test Location 2"},
		{ID: 3, Name: "Test Location 3"},
	}

	t.Run("when service returns locations successfully", func(t *testing.T) {
		s.service.
			On("ListAll", mock.Anything, "", int64(10), repo.Filters{}).
			Return(locations, "", nil)

		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		s.handler.ListAll(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("when service returns locations for filters successfully", func(t *testing.T) {
		s.service.
			On("ListAll", mock.Anything, "", int64(10), repo.Filters{"name:ilike": "Test Location 1"}).
			Return(locations, "", nil)

		r := httptest.NewRequest(http.MethodGet, url, nil)
		r.URL.RawQuery = "query=Test Location 1"
		w := httptest.NewRecorder()

		s.handler.ListAll(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("when service returns error", func(t *testing.T) {
		s.service.
			On("ListAll", mock.Anything, "", int64(1), repo.Filters{}).
			Return([]*bitsb.Location{}, "", apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodGet, url, nil)
		r.URL.RawQuery = "limit=1"
		w := httptest.NewRecorder()

		s.handler.ListAll(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func (s *LocationHandlerTestSuite) TestGetByID() {
	t := s.T()

	location := &bitsb.Location{ID: 1, Name: "Test Location 1"}

	t.Run("when service returns location successfully", func(t *testing.T) {
		s.service.
			On("GetByID", mock.Anything, int64(1)).
			Return(location, nil)

		r := httptest.NewRequest(http.MethodGet, "/location/1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.GetByID(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("when service returns error", func(t *testing.T) {
		s.service.
			On("GetByID", mock.Anything, int64(3)).
			Return(&bitsb.Location{}, apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodGet, "/location/3", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.GetByID(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("when the url param is invalid", func(t *testing.T) {
		s.service.
			On("GetByID", mock.Anything, int64(3)).
			Return(&bitsb.Location{}, apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodGet, "/location/invalid", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.GetByID(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func (s *LocationHandlerTestSuite) TestCreate() {
	t := s.T()
	t.Skip() //TODO: fix tests

	location := &bitsb.Location{Name: "Test Location 1"}

	t.Run("when service returns location successfully", func(t *testing.T) {
		s.service.
			On("Create", mock.Anything, &bitsb.Location{
				Name: "Test Location 1",
			}).
			Return(&bitsb.Location{ID: 1, Name: "Test Location 1"}, nil)

		body, _ := json.Marshal(location)

		r := httptest.NewRequest(http.MethodPost, "/location", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.handler.Create(w, r)
		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("when service returns error", func(t *testing.T) {
		s.service.
			On("Create", mock.Anything, mock.Anything).
			Return(&bitsb.Location{}, apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodPost, "/location", nil)
		w := httptest.NewRecorder()

		s.handler.Create(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("when the body is invalid", func(t *testing.T) {
		s.service.
			On("Create", mock.Anything, mock.Anything).
			Return(&bitsb.Location{}, apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodPost, "/location", strings.NewReader("invalid"))
		w := httptest.NewRecorder()

		s.handler.Create(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func (s *LocationHandlerTestSuite) TestUpdate() {}

func (s *LocationHandlerTestSuite) TestDelete() {
	t := s.T()

	t.Run("when service returns location successfully", func(t *testing.T) {
		s.service.
			On("Delete", mock.Anything, int64(1)).
			Return(nil)

		r := httptest.NewRequest(http.MethodDelete, "/location/1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.Delete(w, r)
		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("when service returns error", func(t *testing.T) {
		s.service.
			On("Delete", mock.Anything, int64(3)).
			Return(apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodDelete, "/location/3", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.Delete(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("when the url param is invalid", func(t *testing.T) {
		s.service.
			On("Delete", mock.Anything, int64(3)).
			Return(apperrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodDelete, "/location/invalid", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		s.handler.Delete(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
