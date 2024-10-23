package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/IceTweak/url-shortener/internal/lib/api/response"
	"github.com/IceTweak/url-shortener/internal/lib/logger/sl"
	"github.com/IceTweak/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	aliasParam = "alias"
)

// URLGetter is an interface for getting url by alias
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	// Defined in storage.sqlite.GetURL
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, aliasParam)
		if alias == "" {
			log.Error("alias parameter is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		log.Info("get redirect with alias", slog.String("alias", alias))

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("url not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info(
			"get url successfully",
			slog.String("url", url),
			slog.String("alias", alias),
		)

		http.Redirect(w, r, url, http.StatusFound)
	}
}
