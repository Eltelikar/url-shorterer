package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "project_1/internal/lib/api/response"
	"project_1/internal/lib/logger/sl"
	"project_1/internal/lib/random"
	"project_1/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 6

// URLSaaver defines the interface for saving URLs.
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

type AliasCheck interface {
	AliasExists(alias string) (bool, error)
}

func New(log *slog.Logger, stgInterface interface {
	URLSaver
	AliasCheck
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//TODO: fix request_id output
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
			check, err := stgInterface.AliasExists(alias)
			if err != nil {
				log.Error("failed to check exists alias", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to check alias exists"))
				return
			}

			if check {
				alias = random.NewRandomString(aliasLength)
			}

		}

		id, err := stgInterface.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url arleady exists", slog.String("key", req.URL))

			render.JSON(w, r, resp.Error("url arleady exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved successfully", slog.Int64("id", id), slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
