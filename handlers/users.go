package handlers

import (
	"encoding/json"
	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common"
	errors2 "github.com/misgorod/co-dev/errors"
	"github.com/misgorod/co-dev/models"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersHandler struct {
	Client *mongo.Client
}

func (u *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, errors2.ErrNoID)
		return
	}

	user, err := models.GetUser(r.Context(), u.Client, id)
	if err != nil {
		common.RespondError(w, err)
		return
	}

	common.RespondJSON(w, http.StatusOK, user)
}

func (u *UsersHandler) Put(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	var info models.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
		return
	}
	user, err := models.PutUser(r.Context(), u.Client, userID, &info)
	if err != nil {
		common.RespondError(w, err)
		return
	}

	common.RespondJSON(w, 200, user)
}

func (u *UsersHandler) PostImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	user, err := models.GetUser(r.Context(), u.Client, userID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	r.ParseMultipartForm(16 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			common.RespondError(w, errors2.ErrNoFileKey)
		}
		common.RespondError(w, err)
		return
	}
	defer file.Close()
	fileObj, err := models.UploadUserImage(r.Context(), u.Client, file, user)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, 201, fileObj)
}
