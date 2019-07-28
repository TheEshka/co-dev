package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/misgorod/co-dev/common"
	uerrors "github.com/misgorod/co-dev/users/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Client *mongo.Client
}

func (a *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, uerrors.ErrNoID)
		return
	}

	user, err := GetUser(r.Context(), a.Client, id)
	if err != nil {
		common.RespondError(w, err)
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}
