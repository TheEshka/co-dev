package post

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
	Client *mongo.Client
}

type pageOptions struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (p *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
	var post *Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	userId, ok := auth.GetUserId(r.Context())
	if !ok {
		common.RespondError(w, http.StatusUnauthorized, "Token is invalid")
	}

	post, err := CreatePost(r.Context(), p.Client, userId, post)
	if err != nil {
		switch err {
		case primitive.ErrInvalidHex:
			common.RespondError(w, http.StatusUnauthorized, "Token is invalid")
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}

	common.RespondJSON(w, http.StatusCreated, post)
}

func (p *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var pageOptions *pageOptions
	if err := json.NewDecoder(r.Body).Decode(&pageOptions); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	if pageOptions.Limit == 0 || pageOptions.Limit > 50 {
		pageOptions.Limit = 50
	}

	posts, err := GetPosts(r.Context(), p.Client, pageOptions.Offset, pageOptions.Limit)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		return
	}

	common.RespondJSON(w, http.StatusOK, posts)
}

func (p *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "ID not specified")
	}

	post, err := GetPost(r.Context(), p.Client, id)
	if err != nil {
		switch err {
		case ErrPostNotFound:
			common.RespondError(w, http.StatusNotFound, err.Error())
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
		return
	}

	common.RespondJSON(w, http.StatusOK, post)
}

func (p *PostHandler) MemberPost(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserId(r.Context())
	if !ok {
		common.RespondError(w, http.StatusUnauthorized, "Token is invalid")
	}
	postId := chi.URLParam(r, "id")
	err := AddMember(r.Context(), p.Client, postId, userId)
	if err != nil {
		switch err {
		case auth.ErrWrongToken:
			common.RespondError(w, http.StatusUnauthorized, err.Error())
		case ErrPostNotFound:
			common.RespondError(w, http.StatusNotFound, err.Error())
		case ErrMemberAlreadyExists:
			common.RespondError(w, http.StatusBadRequest, err.Error())
		case ErrMemberIsAuthor:
			common.RespondError(w, http.StatusConflict, err.Error())
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
	}
}

func (p *PostHandler) MemberDelete(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserId(r.Context())
	if !ok {
		common.RespondError(w, http.StatusUnauthorized, "Token is invalid")
	}
	postId := chi.URLParam(r, "id")
	err := DeleteMember(r.Context(), p.Client, postId, userId)
	if err != nil {
		switch err {
		case auth.ErrWrongToken:
			common.RespondError(w, http.StatusUnauthorized, err.Error())
		case ErrPostNotFound:
			fallthrough
		case ErrMebmerNotExists:
			common.RespondError(w, http.StatusNotFound, err.Error())
		default:
			common.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Internal: %s", err.Error()))
		}
	}
}
