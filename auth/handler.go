package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/tokens"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	Client *mongo.Client
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user *User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	user, err := createUser(r.Context(), a.Client, user.Email, user.Password)
	if err != nil {
		switch err {
		case ErrUserExists:
			common.RespondError(w, http.StatusBadRequest, "User with this email already exists")
			break
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}

	token, err := tokens.CreateToken(user.ID.String())
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusCreated, user)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user *User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	user, err := validateUser(r.Context(), a.Client, user.Email, user.Password)
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

	token, err := tokens.CreateToken(user.ID.Hex())
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	common.RespondJSON(w, http.StatusCreated, user)
}
