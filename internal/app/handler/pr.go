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

type PostPullRequestReassignResponse struct {
	Pr         api.PullRequest `json:"pr"`
	ReplacedBy string          `json:"replaced_by"`
}

type PrHandler struct {
	prService *service.PrService
}

func NewPrHandler(s *service.PrService) *PrHandler {
	return &PrHandler{prService: s}
}

func (h *PrHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var body api.PostPullRequestCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	prID := strings.TrimSpace(body.PullRequestId)
	prName := strings.TrimSpace(body.PullRequestName)
	authorID := strings.TrimSpace(body.AuthorId)

	if prID == "" || prName == "" || authorID == "" {
		http.Error(w, "pull_request_id, pull_request_name and author_id must not be empty", http.StatusBadRequest)
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), prID, prName, authorID)
	if err != nil {
		switch err {
		case domain_errors.ErrPullRequestExists:
			WriteJSONError(w, http.StatusConflict, api.PREXISTS, "pull request already exists")
			return
		case domain_errors.ErrUserNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "author not found")
			return
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := mapper.ToAPIPullRequest(pr)
	WriteJSON(w, http.StatusCreated, map[string]interface{}{"pr": resp})
}

func (h *PrHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var body api.PostPullRequestMergeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	prID := strings.TrimSpace(body.PullRequestId)
	if prID == "" {
		http.Error(w, "pull_request_id must not be empty", http.StatusBadRequest)
		return
	}

	pr, err := h.prService.Merge(r.Context(), prID)
	if err != nil {
		switch err {
		case domain_errors.ErrPullRequestNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "pull request not found")
			return
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := mapper.ToAPIPullRequest(pr)
	WriteJSON(w, http.StatusOK, map[string]interface{}{"pr": resp})
}

func (h *PrHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var body api.PostPullRequestReassignJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	prID := strings.TrimSpace(body.PullRequestId)
	oldReviewerID := strings.TrimSpace(body.OldUserId)
	if prID == "" || oldReviewerID == "" {
		http.Error(w, "pull_request_id and old_user_id must not be empty", http.StatusBadRequest)
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(r.Context(), prID, oldReviewerID)
	if err != nil {
		switch err {
		case domain_errors.ErrPullRequestNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "pull request not found")
			return
		case domain_errors.ErrReviewerNotAssigned:
			WriteJSONError(w, http.StatusConflict, api.NOTASSIGNED, "reviewer not assigned to PR")
			return
		case domain_errors.ErrPRMerged:
			WriteJSONError(w, http.StatusConflict, api.PRMERGED, "cannot reassign merged PR")
			return
		case domain_errors.ErrNoReplacementCandidate:
			WriteJSONError(w, http.StatusConflict, api.NOCANDIDATE, "no replacement candidate available")
			return
		case domain_errors.ErrUserNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "user not found")
			return
		}
	}

	resp := PostPullRequestReassignResponse{
		Pr:         mapper.ToAPIPullRequest(pr),
		ReplacedBy: newReviewerID,
	}

	WriteJSON(w, http.StatusOK, resp)
}
