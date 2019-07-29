package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/common/errors"
	perrors "github.com/misgorod/co-dev/post/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	Client   *mongo.Client
	Validate *validator.Validate
}

type file struct {
	ID primitive.ObjectID `json:"id"`
}

func (p *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var post *Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		common.RespondError(w, errors.ErrDecodeRequest)
		return
	}
	if err := p.Validate.Struct(post); err != nil {
		common.RespondError(w, errors.ErrValidateRequest)
		return
	}
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
	}
	post, err := CreatePost(r.Context(), p.Client, userID, post)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, http.StatusCreated, post)
}

func (p *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	slimit := r.URL.Query().Get("limit)
	if slimit == "" {
		slimit = "10"	
	}
	limit, err := strconv.Atoi(slimit)
	if err != nil {
		common.RespondError(w, errors.ErrDecodeRequest)	
	}
	soffset := r.URL.Query().Get("offset")
	if soffset == "" {
	    soffset = "0"    
	}
	offset, err := strconv.Atoi(soffset)
	if err != nil {
		common.RespondError(w, errors.ErrDecodeRequest)		
	}
	posts, err := GetPosts(r.Context(), p.Client, offset, limit)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, http.StatusOK, posts)
}

func (p *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondError(w, perrors.ErrNoPostID)
	}
	post, err := GetPost(r.Context(), p.Client, id)
	if err != nil {
		common.RespondError(w, err)
		return
	}

	common.RespondJSON(w, http.StatusOK, post)
}

func (p *Handler) MembersPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, perrors.ErrNoPostID)
		return
	}
	err := AddMember(r.Context(), p.Client, postID, userID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *Handler) MemberPut(w http.ResponseWriter, r *http.Request) {
	authorID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, perrors.ErrNoPostID)
		return
	}
	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		common.RespondError(w, perrors.ErrNoMemberID)
		return
	}
	err := ApproveMember(r.Context(), p.Client, postID, authorID, memberID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *Handler) MemberDelete(w http.ResponseWriter, r *http.Request) {
	authorID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, perrors.ErrNoPostID)
		return
	}
	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		common.RespondError(w, perrors.ErrNoMemberID)
		return
	}
	err := DeleteMember(r.Context(), p.Client, postID, authorID, memberID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *Handler) MembersDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, perrors.ErrNoPostID)
		return
	}
	err := DeleteMemberSelf(r.Context(), p.Client, postID, userID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}

func (p *Handler) PostImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		common.RespondError(w, errors.ErrWrongToken)
		return
	}
	postID := chi.URLParam(r, "id")
	if postID == "" {
		common.RespondError(w, perrors.ErrNoPostID)
		return
	}
	post, err := GetPost(r.Context(), p.Client, postID)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	if post.Author.ID.Hex() != userID {
		common.RespondError(w, perrors.ErrNotAnAuthor)
		return
	}
	r.ParseMultipartForm(16 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			common.RespondError(w, errors.ErrNoFileKey)
		}
		common.RespondError(w, err)
		return
	}
	defer file.Close()
	fileObj, err := UploadImage(r.Context(), p.Client, file, post)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, 201, fileObj)
}

func (p *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	err := DownloadImage(r.Context(), p.Client, imageID, w)
	if err != nil {
		common.RespondError(w, err)
		return
	}
}
