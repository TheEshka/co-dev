package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/users"

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
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			common.RespondError(w, http.StatusBadRequest, "Invalid json")
			break
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}
	err := createUser(r.Context(), a.Client, &user)
	if err != nil {
		switch err {
		case users.ErrUserExists:
			common.RespondError(w, http.StatusBadRequest, "User with this email already exists")
			break
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}
	token, err := CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusCreated, &user)
}

func (a *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user loginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}
	if err := a.Validate.Struct(user); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			common.RespondError(w, http.StatusBadRequest, "Invalid json")
			break
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}
	err := validateUser(r.Context(), a.Client, &user)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			common.RespondError(w, http.StatusNotFound, "User not found")
			break
		case bcrypt.ErrMismatchedHashAndPassword:
			common.RespondError(w, http.StatusForbidden, "Wrong password")
			break
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}

	token, err := CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusOK, &user)
}
