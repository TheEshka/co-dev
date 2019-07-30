package handlers

import (
	"github.com/go-chi/chi"
	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

type ImagesHandler struct {
	Client   *mongo.Client
	Validate *validator.Validate
}

func (i *ImagesHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	err := models.DownloadFile(r.Context(), i.Client, imageID, w)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}
