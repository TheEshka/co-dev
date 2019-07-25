package users

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/misgorod/co-dev/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersHandler struct {
	Client *mongo.Client
}

func (a *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "ID not specified")
		return
	}

	user, err := GetUser(r.Context(), a.Client, id)
	if err != nil {
		switch err {
		case ErrUserNotExists:
			common.RespondError(w, http.StatusBadRequest, err.Error())
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}
