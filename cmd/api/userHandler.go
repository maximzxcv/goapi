package main

import (
	"goapi/bal"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUsers .....
func (server *Server) GetUsers(ctx *gin.Context) {
	us := make([]bal.User, 5)
	us[0] = bal.User{
		ID:   uuid.New().String(),
		Name: "Name 1",
	}
	us[1] = bal.User{
		ID:   uuid.New().String(),
		Name: "Name 2",
	}

	ctx.JSON(http.StatusOK, us)
}

// CreateUser ....
func (server *Server) CreateUser(ctx *gin.Context) {
	var req bal.CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := bal.CreateUser{
		Name:            req.Name,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		//	Balance:  0,
	}

	//save to DB
	// arg := db.CreateAccountParams{
	// 	Owner:    req.Owner,
	// 	Currency: req.Currency,
	// 	Balance:  0,
	// }

	// account, err := server.store.CreateAccount(ctx, arg)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	//ctx.JSON(http.StatusOK, account)

	ctx.JSON(http.StatusOK, arg)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
