package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "project_1/internal/lib/api/response"
	"project_1/internal/lib/logger/sl"
	"project_1/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type DeleteURL interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleteURL DeleteURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias is empry"))
			return
		}
		log.Info("alias get succsessfully", slog.String("alias", alias))

		err := deleteURL.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		} else if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url was deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
