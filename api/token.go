package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(c *gin.Context) {

	var req renewAccessTokenRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.store.GetSession(c, refreshPayload.ID)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("session is blocked")))
		return
	}

	if session.Username != refreshPayload.Username {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("incorrect session user")))
		return
	}

	if req.RefreshToken != session.RefreshToken {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("missmatched refresh token")))
		return
	}

	if time.Now().After(refreshPayload.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("expired session")))
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(refreshPayload.Username, time.Duration(s.config.AccesTokenDuration*int(time.Minute)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt,
	})

}
