package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/ryanMiranda98/simplebank/db/sqlc"
	"github.com/ryanMiranda98/simplebank/token"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func ValidateSession(ctx *gin.Context, session *db.Session, request *RenewAccessTokenRequest, payload *token.Payload) error {
	if session.IsBlocked {
		return errors.New("session is blocked")
	}
	if session.RefreshToken != request.RefreshToken {
		return errors.New("session/token mismatch")
	}
	if session.Username != payload.Username {
		return errors.New("incorrect session user")
	}
	if time.Now().After(session.ExpiresAt) {
		return errors.New("session has expired")
	}
	return nil
}

func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	refreshTokenPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = ValidateSession(ctx, &session, &req, refreshTokenPayload); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := RenewAccessTokenResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, resp)
}
