package handlers

import (
	"encoding/json"
	"github.com/misgorod/co-dev/common"
	errors2 "github.com/misgorod/co-dev/errors"
	"github.com/misgorod/co-dev/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/misgorod/co-dev/auth"
	"go.mongodb.org/mongo-driver/mongo"

	"gopkg.in/go-playground/validator.v9"
)

type PostsHandler struct {
	Client   *mongo.Client
	Validate *validator.Validate
}

func (p *PostsHandler) Post(w http.ResponseWriter, r *http.Request) {
	var post *models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
		return
	}
	if err := p.Validate.Struct(post); err != nil {
		common.RespondError(w, errors2.ErrValidateRequest)
		return
	}
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
	}
	post, err := models.CreatePost(r.Context(), p.Client, userID, post)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, http.StatusCreated, post)
}

func (p *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	slimit := r.URL.Query().Get("limit")
	if slimit == "" {
		slimit = "10"
	}
	limit, err := strconv.Atoi(slimit)
	if err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
	}
	soffset := r.URL.Query().Get("offset")
	if soffset == "" {
		soffset = "0"
	}
	offset, err := strconv.Atoi(soffset)
	if err != nil {
		common.RespondError(w, errors2.ErrDecodeRequest)
	}
	posts, err := models.GetPosts(r.Context(), p.Client, offset, limit)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, http.StatusOK, posts)
}

func (p *PostsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, errors2.ErrNoPostID)
	}
	post, err := models.GetPost(r.Context(), p.Client, id)
	if err != nil {
		common.RespondError(w, err)
		return
	}

	common.RespondJSON(w, http.StatusOK, post)
}

func (p *PostsHandler) MembersPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, errors2.ErrNoPostID)
		return
	}
	err := models.AddPostMember(r.Context(), p.Client, postID, userID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *PostsHandler) MemberPut(w http.ResponseWriter, r *http.Request) {
	authorID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, errors2.ErrNoPostID)
		return
	}
	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		common.RespondError(w, errors2.ErrNoMemberID)
		return
	}
	err := models.ApprovePostMember(r.Context(), p.Client, postID, authorID, memberID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *PostsHandler) MemberDelete(w http.ResponseWriter, r *http.Request) {
	authorID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, errors2.ErrNoPostID)
		return
	}
	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		common.RespondError(w, errors2.ErrNoMemberID)
		return
	}
	err := models.DeletePostMember(r.Context(), p.Client, postID, authorID, memberID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *PostsHandler) MembersDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, errors2.ErrNoPostID)
		return
	}
	err := models.DeletePostMemberSelf(r.Context(), p.Client, postID, userID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *PostsHandler) PostImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors2.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, errors2.ErrNoPostID)
		return
	}
	post, err := models.GetPost(r.Context(), p.Client, postID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	if post.Author.ID.Hex() != userID {
		common.RespondError(w, errors2.ErrNotAnAuthor)
		return
	}
	err = r.ParseMultipartForm(16 << 20)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			common.RespondError(w, errors2.ErrNoFileKey)
		}
		common.RespondError(w, err)
		return
	}
	defer file.Close()
	fileObj, err := models.UploadPostImage(r.Context(), p.Client, file, post)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, 201, fileObj)
}
