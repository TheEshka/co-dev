package handlers

import (
	"encoding/json"
	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common"
	errors2 "github.com/misgorod/co-dev/errors"
	"github.com/misgorod/co-dev/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

type AuthHandler struct {
	Client   *mongo.Client
	Validate *validator.Validate
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.RegUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		common.RespondError(w, errors2.ErrValidateRequest)
		return
	}
	err := models.CreateUser(r.Context(), a.Client, &user)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	token, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, err)
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusCreated, &user)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		common.RespondError(w, errors2.ErrValidateRequest)
		return
	}
	err := models.ValidateUser(r.Context(), a.Client, &user)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	token, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, err)
		return
	}
	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusOK, &user)
}
