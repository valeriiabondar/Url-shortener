package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"urlShortener/internal/http-server/handlers/mocks"
	"urlShortener/internal/storage"
	"urlShortener/internal/utils/logger"
)

func TestDeleteUrl(t *testing.T) {
	tests := []struct {
		name           string
		alias          string
		expectedStatus int
		mockError      error
		respError      string
	}{
		{
			name:           "Success",
			alias:          "test_alias",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Empty alias",
			alias:          "",
			expectedStatus: http.StatusBadRequest,
			respError:      "invalid request",
		},
		{
			name:           "Url not found",
			alias:          "non_existing_alias",
			expectedStatus: http.StatusNotFound,
			mockError:      storage.ErrUrlNotFound,
			respError:      "url not found",
		},
		{
			name:           "Internal error",
			alias:          "test_alias",
			mockError:      errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
			respError:      "internal error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewUrlDeleter(t)

			if tc.alias != "" {
				if tc.mockError == nil {
					urlDeleterMock.On("DeleteUrl", tc.alias).Return(nil).Once()
				} else {
					urlDeleterMock.On("DeleteUrl", tc.alias).Return(tc.mockError).Once()
				}
			}

			handler := DeleteUrl(logger.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.alias), nil)
			require.NoError(t, err)

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("alias", tc.alias)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus != http.StatusNoContent {
				body := rr.Body.String()
				var resp Response
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
