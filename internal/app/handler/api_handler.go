package handler

import (
	"encoding/json"
	"net/http"
	"test/internal/api"
)

type APIHandler struct {
	team *TeamHandler
	user *UserHandler
	pr   *PrHandler
}

func NewAPIHandler(team *TeamHandler, user *UserHandler, pr *PrHandler) *APIHandler {
	return &APIHandler{
		team: team,
		user: user,
		pr:   pr,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func WriteJSONError(w http.ResponseWriter, status int, code api.ErrorResponseErrorCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := api.ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = msg

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *APIHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	h.pr.PostPullRequestCreate(w, r)
}

func (h *APIHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	h.pr.PostPullRequestMerge(w, r)
}

func (h *APIHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	h.pr.PostPullRequestReassign(w, r)
}

func (h *APIHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	h.team.PostTeamAdd(w, r)
}

func (h *APIHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	h.team.GetTeamGet(w, r, params)
}

func (h *APIHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	h.user.GetUsersGetReview(w, r, params)
}

func (h *APIHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	h.user.PostUsersSetIsActive(w, r)
}
