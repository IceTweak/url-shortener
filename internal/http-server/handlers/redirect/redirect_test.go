package redirect_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IceTweak/url-shortener/internal/http-server/handlers/redirect"
	"github.com/IceTweak/url-shortener/internal/http-server/handlers/redirect/mocks"
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

func TestRedirectHandler(t *testing.T) {
	testcases := []TestCase{}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respErr == "" || tc.mockErr != nil {
				urlGetterMock.On("GetURL", tc.alias, mock.AnythingOfType("string")).
					Return("", tc.mockErr).
					Once()
			}

			// DiscardLogger - logger that do nothing (logger mock)
			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

			path := fmt.Sprintf("/%s", tc.alias)

			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			redirectedURL := rr.Result().Header.Get("Location")

			require.Equal(t, rr.Code, http.StatusFound)
			require.Equal(t, redirectedURL, tc.url)
		})
	}
}
