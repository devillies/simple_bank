package api

import (
	"net/http/httptest"
	"testing"

	mockdb "github.com/devillies/simple_bank/db/mock"
	"github.com/gin-gonic/gin"
)

func TestCreateUserApi(t *testing.T) {

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{}
}
