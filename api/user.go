package api

import (
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

	ctx.IndentedJSON(http.StatusAccepted, user)

}
