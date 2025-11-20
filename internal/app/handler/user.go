package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"test/internal/api"
	"test/internal/app/mapper"
	"test/internal/domain/domain_errors"
	"test/internal/domain/service"
)

type UsersGetReviewResponse struct {
	UserID       string                 `json:"user_id"`
	PullRequests []api.PullRequestShort `json:"pull_requests"`
}

type UserHandler struct {
	userService *service.UserService
	prService   *service.PrService
}

func NewUserHandler(us *service.UserService, ps *service.PrService) *UserHandler {
	return &UserHandler{userService: us, prService: ps}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	userId := strings.TrimSpace(params.UserId)
	if userId == "" {
		http.Error(w, "user_id must not be empty", http.StatusBadRequest)
		return
	}

	prList, err := h.prService.GetByReviewer(r.Context(), userId)
	if err != nil {
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UsersGetReviewResponse{
		UserID:       userId,
		PullRequests: make([]api.PullRequestShort, 0, len(prList)),
	}

	for _, pr := range prList {
		resp.PullRequests = append(resp.PullRequests, mapper.ToAPIPullRequestShort(pr))
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var body api.PostUsersSetIsActiveJSONBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.UserId) == "" {
		http.Error(w, "user_id must not be empty", http.StatusBadRequest)
		return
	}

	u, err := h.userService.SetIsActive(r.Context(), body.UserId, body.IsActive)
	if err != nil {
		switch err {
		case domain_errors.ErrUserNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "user not found")
			return
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	apiUser := mapper.ToAPIUser(u)

	WriteJSON(w, http.StatusOK, map[string]interface{}{"team": apiUser})
}
