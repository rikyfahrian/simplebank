package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "techschool/db/mock"
	db "techschool/db/sqlc"
	"techschool/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

	account := randAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	//build stubs
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	//start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req)

	//check response
	require.Equal(t, http.StatusOK, recorder.Code)

}

func randAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomName(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}
