package api

import (
	"database/sql"
	"net/http"

	db "github.com/devillies/simple_bank/db/sqlc"
	"github.com/devillies/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6" `
	FullName string `json:"fullname" binding:"required" `
	Email    string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	password, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: password,
		FullName:     req.FullName,
		Email:        req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.IndentedJSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := convertToUserResponse(user)

	ctx.IndentedJSON(http.StatusAccepted, res)

}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {

			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := convertToUserResponse(user)

	ctx.IndentedJSON(http.StatusOK, response)

}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUser(ctx, arg)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []db.UserResponse

	for _, user := range users {
		convertUser := convertToUserResponse(user)
		response = append(response, convertUser)
	}
	ctx.IndentedJSON(http.StatusOK, response)
}

func convertToUserResponse(usr db.User) db.UserResponse {
	return db.UserResponse{
		Username:          usr.Username,
		FullName:          usr.FullName,
		Email:             usr.Email,
		PasswordChangedAt: usr.PasswordChangedAt,
		CreatedAt:         usr.CreatedAt,
	}
}
