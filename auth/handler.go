package auth

import (
	"encoding/json"
	"net/http"

	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/common/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	Client   *mongo.Client
	Validate *validator.Validate
}

func (a *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user regUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, errors.ErrDecodeRequest)
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		common.RespondError(w, errors.ErrValidateRequest)
		return
	}
	err := createUser(r.Context(), a.Client, &user)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	token, err := CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, err)
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusCreated, &user)
}

func (a *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user loginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, errors.ErrDecodeRequest)
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		common.RespondError(w, errors.ErrValidateRequest)
		return
	}
	err := validateUser(r.Context(), a.Client, &user)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	token, err := CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, err)
		return
	}
	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusOK, &user)
}
