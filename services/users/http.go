package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fayleenpc/tj-jeans/internal/auth"
	"github.com/fayleenpc/tj-jeans/internal/types"
	"github.com/fayleenpc/tj-jeans/internal/utils"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlerClient struct {
	client pb.UserServiceClient
}

func NewHandlerClient(client pb.UserServiceClient) *HandlerClient {
	return &HandlerClient{client: client}
}

type HandlerHTTP struct {
	client types.UserService
}

func NewHandlerHTTP(client types.UserService) *HandlerHTTP {
	return &HandlerHTTP{client: client}
}

func (h *HandlerHTTP) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", h.handleGetUsers_Proto)
	mux.HandleFunc("GET /api/v1/users/{user_id}", h.handleGetUserByID_Proto)
	mux.HandleFunc("PATCH /api/v1/users/{user_id}/update", h.handleUpdateUserByID_Proto)
	mux.HandleFunc("DELETE /api/v1/users/{user_id}/delete", h.handleDeleteUserByID_Proto)
}

func (h *HandlerHTTP) handleGetUsers_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetUsers_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())

	if userRole == "admin" {
		response, err := h.client.GetUsers(r.Context(), &pb.GetUsersRequest{})
		rStatus := status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.ResourceExhausted {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}
		utils.WriteJSON(w, http.StatusOK, response.GetUsers())
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission, denied"))
	}
}

func (h *HandlerHTTP) handleGetUserByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleGetUserByID_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		response, err := h.client.GetUserByID(r.Context(), &pb.GetUserByIDRequest{Id: int32(userID)})
		rStatus := status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}

		utils.WriteJSON(w, http.StatusOK, response.GetUser())
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission, denied"))
	}
}

func (h *HandlerHTTP) handleUpdateUserByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleUpdateUserByID_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var payload pb.User
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		payload.Id = int32(userID)

		oldUser, err := h.client.GetUserByID(r.Context(), &pb.GetUserByIDRequest{Id: int32(userID)})
		rStatus := status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}

		response, err := h.client.UpdateUser(r.Context(), &pb.UpdateUserRequest{User: &payload})
		rStatus = status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}
		updatedUser, err := h.client.GetUserByID(r.Context(), &pb.GetUserByIDRequest{Id: int32(response.GetUpdatedCount())})
		rStatus = status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"updated_id": response.GetUpdatedCount(), "old_user": oldUser.GetUser(), "updated_user": updatedUser.GetUser()})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission, denied"))
	}
}

func (h *HandlerHTTP) handleDeleteUserByID_Proto(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("handleDeleteUserByID_Proto")
	defer span.Finish()

	span.SetTag(string(ext.Component), "http")
	span.SetTag("http.method", r.Method)
	userRole := auth.GetUserRoleFromContext(r.Context())
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if userRole == "admin" {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		oldUserResponse, err := h.client.GetUserByID(r.Context(), &pb.GetUserByIDRequest{Id: int32(userID)})
		rStatus := status.Convert(err)
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}
		deletedUserResponse, err := h.client.DeleteUserByID(r.Context(), &pb.DeleteUserByIDRequest{Id: oldUserResponse.GetUser().GetId()})
		if rStatus != nil {
			if rStatus.Code() != codes.InvalidArgument {
				utils.WriteError(w, http.StatusBadRequest, rStatus.Err())
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Err())
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"deleted_id": deletedUserResponse.GetDeletedCount(), "deleted_user": oldUserResponse.GetUser()})
	} else {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("permission, denied"))
	}

}
