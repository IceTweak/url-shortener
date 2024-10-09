package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IceTweak/url-shortener/internal/http-server/handlers/url/save"
	"github.com/IceTweak/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/IceTweak/url-shortener/internal/lib/logger/sl/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name    string
	alias   string
	url     string
	respErr string
	mockErr error
}

func TestSaveHandler(t *testing.T) {
	testcases := []TestCase{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:    "Empty URL",
			url:     "",
			alias:   "some_alias",
			respErr: "field URL is a required field",
		},
		{
			name:    "Invalid URL",
			url:     "some invalid URL",
			alias:   "some_alias",
			respErr: "field URL is not a valid URL",
		},
		{
			name:    "SaveURL Error",
			alias:   "test_alias",
			url:     "https://google.com",
			respErr: "failed to add url",
			mockErr: errors.New("unexpected error"),
		},
	}

	for _, tc := range testcases {
		//? Parallel running issue, refs:
		// https://qna.habr.com/q/814247
		// https://go.dev/play/p/tEl60VkEtQv
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			//? Test parallel running ref:
			// https://engineering.mercari.com/en/blog/entry/20220408-how_to_use_t_parallel/
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respErr == "" || tc.mockErr != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockErr).
					Once()
			}

			// DiscardLogger - logger that do nothing (logger mock)
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"alias": "%s", "url": "%s"}`, tc.alias, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respErr, resp.Error)

			// TODO: add more checks
		})
	}
}
