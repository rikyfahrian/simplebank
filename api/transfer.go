package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "techschool/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req transferRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, isOK := s.validAccount(c, req.FromAccountID, req.Currency)
	if !isOK {
		return
	}

	_, isOK = s.validAccount(c, req.ToAccountID, req.Currency)
	if !isOK {
		return
	}

	if fromAccount.Balance < req.Amount {

		c.JSON(http.StatusBadRequest, errorResponse(errors.New("balance is not enough to transfers")))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)

}

func (s *Server) validAccount(c *gin.Context, accountID int64, crncy string) (db.Account, bool) {

	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != crncy {
		err := fmt.Errorf("error different currency %s vs %s", crncy, account.Currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true

}
