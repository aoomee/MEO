package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"miaomiaowux/internal/auth"
	"miaomiaowux/internal/storage"
)

type PackageAssignmentsHandler struct {
	repo         *storage.TrafficRepository
	remoteManage *RemoteManageHandler
	pusher       *LimiterConfigPusher
	admin        bool
}

func NewAdminPackageAssignmentsHandler(repo *storage.TrafficRepository, remoteManage *RemoteManageHandler, pusher *LimiterConfigPusher) http.Handler {
	return &PackageAssignmentsHandler{repo: repo, remoteManage: remoteManage, pusher: pusher, admin: true}
}

func NewUserPackageAssignmentsHandler(repo *storage.TrafficRepository) http.Handler {
	return &PackageAssignmentsHandler{repo: repo}
}

type packageAssignmentPayload struct {
	AssignmentID   int64    `json:"assignment_id"`
	Username       string   `json:"username"`
	PackageID      int64    `json:"package_id"`
	StartDate      string   `json:"start_date"`
	ExpireDate     string   `json:"expire_date"`
	IsReset        *bool    `json:"is_reset"`
	ResetDay       *int     `json:"reset_day"`
	TrafficLimitGB *float64 `json:"traffic_limit_override_gb"`
}

type packageAssignmentResponse struct {
	storage.UserPackageAssignment
	UsedUplink   int64                                   `json:"used_uplink"`
	UsedDownlink int64                                   `json:"used_downlink"`
	UsedTotal    int64                                   `json:"used_total"`
	DailyTraffic []storage.PackageAssignmentDailyTraffic `json:"daily_traffic,omitempty"`
}

func parseAssignmentDates(pkg *storage.Package, startText, endText string) (time.Time, time.Time, error) {
	start := time.Now()
	var err error
	if strings.TrimSpace(startText) != "" {
		start, err = time.Parse("2006-01-02", startText)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("start_date must use YYYY-MM-DD")
		}
	}
	end := start.AddDate(0, 1, 0)
	if pkg != nil && pkg.CycleDays > 0 {
		end = start.AddDate(0, 0, pkg.CycleDays)
	}
	if strings.TrimSpace(endText) != "" {
		end, err = time.Parse("2006-01-02", endText)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("expire_date must use YYYY-MM-DD")
		}
	}
	if !end.After(start) {
		return time.Time{}, time.Time{}, errors.New("expire_date must be after start_date")
	}
	return start, end, nil
}

func (h *PackageAssignmentsHandler) assignmentResponses(ctx context.Context, assignments []storage.UserPackageAssignment, includeDaily bool) []packageAssignmentResponse {
	out := make([]packageAssignmentResponse, 0, len(assignments))
	end := time.Now().Format("2006-01-02")
	start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	for _, assignment := range assignments {
		up, down, _ := h.repo.GetPackageAssignmentBillableTrafficByDirection(ctx, assignment.ID)
		item := packageAssignmentResponse{UserPackageAssignment: assignment, UsedUplink: up, UsedDownlink: down, UsedTotal: up + down}
		if includeDaily {
			item.DailyTraffic, _ = h.repo.ListPackageAssignmentDailyTraffic(ctx, assignment.ID, start, end)
		}
		out = append(out, item)
	}
	return out
}

func (h *PackageAssignmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	if !h.admin {
		username = auth.UsernameFromContext(r.Context())
	}
	if r.Method == http.MethodGet {
		if username == "" {
			writeError(w, http.StatusBadRequest, errors.New("username is required"))
			return
		}
		assignments, err := h.repo.ListUserPackageAssignments(r.Context(), username, h.admin && r.URL.Query().Get("history") == "1")
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"assignments": h.assignmentResponses(r.Context(), assignments, true)})
		return
	}
	if !h.admin {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	var req packageAssignmentPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	if req.Username == "" {
		req.Username = username
	}
	switch r.Method {
	case http.MethodPost:
		if req.Username == "" || req.PackageID <= 0 {
			writeError(w, http.StatusBadRequest, errors.New("username and package_id are required"))
			return
		}
		pkg, err := h.repo.GetPackage(r.Context(), req.PackageID)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		start, end, err := parseAssignmentDates(pkg, req.StartDate, req.ExpireDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		isReset, resetDay := pkg.IsReset, pkg.ResetDay
		if req.IsReset != nil {
			isReset = *req.IsReset
		}
		if req.ResetDay != nil {
			resetDay = *req.ResetDay
		}
		if isReset && (resetDay < 1 || resetDay > 31) {
			writeError(w, http.StatusBadRequest, errors.New("reset_day must be between 1 and 31"))
			return
		}
		if req.TrafficLimitGB != nil && *req.TrafficLimitGB < 0 {
			writeError(w, http.StatusBadRequest, errors.New("traffic_limit_override_gb cannot be negative"))
			return
		}
		assignment, err := h.addAndProvision(r.Context(), req.Username, req.PackageID, start, end, isReset, resetDay, req.TrafficLimitGB)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"assignment": h.assignmentResponses(r.Context(), []storage.UserPackageAssignment{*assignment}, true)[0]})
	case http.MethodPut:
		if req.AssignmentID <= 0 {
			writeError(w, http.StatusBadRequest, errors.New("assignment_id is required"))
			return
		}
		a, err := h.repo.GetUserPackageAssignment(r.Context(), req.AssignmentID)
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		if req.Username != "" && req.Username != a.Username {
			writeError(w, http.StatusForbidden, errors.New("assignment owner mismatch"))
			return
		}
		pkg, _ := h.repo.GetPackage(r.Context(), a.PackageID)
		if req.StartDate != "" || req.ExpireDate != "" {
			startText, endText := req.StartDate, req.ExpireDate
			if startText == "" && a.PackageStartDate != nil {
				startText = a.PackageStartDate.Format("2006-01-02")
			}
			if endText == "" && a.PackageEndDate != nil {
				endText = a.PackageEndDate.Format("2006-01-02")
			}
			start, end, parseErr := parseAssignmentDates(pkg, startText, endText)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, parseErr)
				return
			}
			a.PackageStartDate, a.PackageEndDate = &start, &end
		}
		if req.IsReset != nil {
			a.IsReset = *req.IsReset
		}
		if req.ResetDay != nil {
			a.ResetDay = *req.ResetDay
		}
		if a.IsReset && (a.ResetDay < 1 || a.ResetDay > 31) {
			writeError(w, http.StatusBadRequest, errors.New("reset_day must be between 1 and 31"))
			return
		}
		if req.TrafficLimitGB != nil {
			if *req.TrafficLimitGB < 0 {
				writeError(w, http.StatusBadRequest, errors.New("traffic_limit_override_gb cannot be negative"))
				return
			}
			bytes := int64(*req.TrafficLimitGB * 1024 * 1024 * 1024)
			a.TrafficLimitOverride = &bytes
		}
		if err := h.repo.UpdateUserPackageAssignment(r.Context(), *a); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"assignment": h.assignmentResponses(r.Context(), []storage.UserPackageAssignment{*a}, true)[0]})
	case http.MethodDelete:
		if req.AssignmentID <= 0 {
			if id, _ := strconv.ParseInt(r.URL.Query().Get("assignment_id"), 10, 64); id > 0 {
				req.AssignmentID = id
			}
		}
		if req.AssignmentID <= 0 {
			writeError(w, http.StatusBadRequest, errors.New("assignment_id is required"))
			return
		}
		if err := h.removeAndDeprovision(r.Context(), req.AssignmentID, req.Username); err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "package assignment removed"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (h *PackageAssignmentsHandler) addAndProvision(ctx context.Context, username string, packageID int64, start, end time.Time, isReset bool, resetDay int, trafficLimitGB *float64) (*storage.UserPackageAssignment, error) {
	assignment, err := h.repo.CreateUserPackageAssignment(ctx, username, packageID, start, end, isReset, resetDay)
	if err != nil {
		return nil, err
	}
	if trafficLimitGB != nil {
		bytes := int64(*trafficLimitGB * 1024 * 1024 * 1024)
		assignment.TrafficLimitOverride = &bytes
		if err := h.repo.UpdateUserPackageAssignment(ctx, *assignment); err != nil {
			_ = h.repo.DeleteUserPackageAssignment(context.WithoutCancel(ctx), assignment.ID)
			return nil, err
		}
	}
	base, err := h.repo.GetUser(ctx, username)
	if err != nil {
		_ = h.repo.DeleteUserPackageAssignment(context.WithoutCancel(ctx), assignment.ID)
		return nil, err
	}
	pkg, err := h.repo.GetPackage(ctx, packageID)
	if err != nil {
		_ = h.repo.DeleteUserPackageAssignment(context.WithoutCancel(ctx), assignment.ID)
		return nil, err
	}
	effectiveNodes, err := effectivePackageNodeIDs(ctx, h.repo, pkg.Nodes, pkg.NodesConfigured)
	if err != nil {
		_ = h.repo.DeleteUserPackageAssignment(context.WithoutCancel(ctx), assignment.ID)
		return nil, err
	}
	user := userForPackageAssignment(base, *assignment)
	updater := &PackageUpdateHandler{repo: h.repo, remoteManage: h.remoteManage, pusher: h.pusher}
	if err := updater.syncPackageUserNodesTransactionally(ctx, []storage.User{user}, nil, effectiveNodes); err != nil {
		_ = h.repo.DeleteUserPackageAssignment(context.WithoutCancel(ctx), assignment.ID)
		return nil, fmt.Errorf("套餐实例下发失败: %w", err)
	}
	provisionWGLeasesForNodes(ctx, h.repo, h.remoteManage, user, effectiveNodes)
	if h.pusher != nil {
		go h.pusher.PushToAllServersForUser(context.Background(), username)
	}
	return assignment, nil
}

func (h *PackageAssignmentsHandler) removeAndDeprovision(ctx context.Context, assignmentID int64, expectedUsername string) error {
	assignment, err := h.repo.GetUserPackageAssignment(ctx, assignmentID)
	if err != nil {
		return err
	}
	if expectedUsername != "" && expectedUsername != assignment.Username {
		return errors.New("assignment owner mismatch")
	}
	base, err := h.repo.GetUser(ctx, assignment.Username)
	if err != nil {
		return err
	}
	pkg, err := h.repo.GetPackage(ctx, assignment.PackageID)
	if err != nil {
		return err
	}
	effectiveNodes, err := effectivePackageNodeIDs(ctx, h.repo, pkg.Nodes, pkg.NodesConfigured)
	if err != nil {
		return err
	}
	user := userForPackageAssignment(base, *assignment)
	updater := &PackageUpdateHandler{repo: h.repo, remoteManage: h.remoteManage, pusher: h.pusher}
	if err := updater.syncPackageUserNodesTransactionally(ctx, []storage.User{user}, effectiveNodes, nil); err != nil {
		return fmt.Errorf("套餐实例解绑失败: %w", err)
	}
	releaseWGLeasesByAssignment(ctx, h.repo, h.remoteManage, assignmentID)
	if err := h.repo.DeleteUserPackageAssignment(ctx, assignmentID); err != nil {
		return err
	}
	if h.pusher != nil {
		go h.pusher.PushToAllServersForUser(context.Background(), assignment.Username)
	}
	return nil
}
