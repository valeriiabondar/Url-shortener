package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"urlShortener/internal/http-server/handlers/mocks"
	"urlShortener/internal/storage"
	"urlShortener/internal/utils/logger"
)

func TestSaveUrl(t *testing.T) {
	tests := []struct {
		name        string
		alias       string
		url         string
		respError   string
		mockError   error
		aliasExists bool
		isGenerated bool
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:        "Empty existing alias",
			alias:       "",
			url:         "https://google.com",
			respError:   "alias already exists",
			mockError:   storage.ErrUrlExists,
			aliasExists: true,
			isGenerated: true,
		},
		{
			name:        "Empty non existing alias",
			alias:       "",
			url:         "https://google.com",
			isGenerated: true,
		},
		{
			name:        "Existing alias",
			alias:       "existing_alias",
			url:         "https://google.com",
			respError:   "alias already exists",
			mockError:   storage.ErrUrlExists,
			aliasExists: true,
		},
		{
			name:  "Non existing alias",
			alias: "non_existing_alias",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "validation error: [field Url is a required field]",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "validation error: [field Url is not a valid URL]",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "could not save url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewUrlSaver(t)

			if tc.isGenerated && tc.aliasExists {
				urlSaverMock.On("AliasExists", mock.AnythingOfType("string")).
					Return(true, nil).Once()
				urlSaverMock.On("AliasExists", mock.AnythingOfType("string")).
					Return(false, nil).Once()
			} else if tc.isGenerated {
				urlSaverMock.On("AliasExists", mock.AnythingOfType("string")).
					Return(tc.aliasExists, nil).Once()
			}

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveUrl", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).Once()
			}

			handler := SaveUrl(logger.NewDiscardLogger(), urlSaverMock)

			reqBody := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(reqBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
