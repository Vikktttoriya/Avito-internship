package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"test/internal/api"
	"test/internal/app/mapper"
	"test/internal/domain/domain_errors"
	"test/internal/domain/model"
	"test/internal/domain/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var body api.Team

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	teamName := strings.TrimSpace(body.TeamName)
	if teamName == "" {
		http.Error(w, "team_name must not be empty", http.StatusBadRequest)
		return
	}

	members := make([]*model.User, 0, len(body.Members))

	for i, m := range body.Members {
		uid := strings.TrimSpace(m.UserId)
		username := strings.TrimSpace(m.Username)

		if uid == "" || username == "" {
			http.Error(w,
				"user_id and username must not be empty for member index "+string(rune(i)),
				http.StatusBadRequest)
			return
		}

		members = append(members,
			model.NewUser(uid, username, teamName, m.IsActive),
		)
	}

	team, err := h.teamService.CreateTeam(r.Context(), teamName, members)
	if err != nil {
		switch err {
		case domain_errors.ErrTeamExists:
			WriteJSONError(w, http.StatusBadRequest, api.TEAMEXISTS, "team already exists")
			return
		case domain_errors.ErrUserHasOpenPullRequests:
			WriteJSONError(w, http.StatusBadRequest, api.PREXISTS, "some users have open PRs")
			return
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := mapper.ToAPITeam(team)
	WriteJSON(w, http.StatusCreated, map[string]interface{}{"team": resp})
}

func (h *TeamHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	teamName := strings.TrimSpace(params.TeamName)
	if teamName == "" {
		http.Error(w, "team_name must not be empty", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), params.TeamName)
	if err != nil {
		switch err {
		case domain_errors.ErrTeamNotFound:
			WriteJSONError(w, http.StatusNotFound, api.NOTFOUND, "team not found")
			return
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := mapper.ToAPITeam(team)
	WriteJSON(w, http.StatusOK, map[string]interface{}{"team": resp})
}
