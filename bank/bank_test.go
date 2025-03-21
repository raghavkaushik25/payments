package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserInfo(t *testing.T) {
	um := NewUserManager()
	oI := OnboartingInfo()
	ru := um.OnboardUsers(oI)
	b := &bank{
		registeredUsers: ru,
	}
	type tests struct {
		name              string
		user              string
		expectedUserName  string
		expectedAccountId string
	}
	tcs := []tests{
		{
			name:              "getAdam",
			user:              "Adam",
			expectedUserName:  "Adam",
			expectedAccountId: b.registeredUsers["Adam"].accountInfo.accountId,
		},
		{
			name:              "getJane",
			user:              "Jane",
			expectedUserName:  "Jane",
			expectedAccountId: b.registeredUsers["Jane"].accountInfo.accountId,
		},
		{
			name:              "getMark",
			user:              "Mark",
			expectedUserName:  "Mark",
			expectedAccountId: b.registeredUsers["Mark"].accountInfo.accountId,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/userInfo?userName=%v", tc.user), nil)
			//req.URL.Query().Set("userName", tc.user)
			req.Header.Set("Content-Type", "application/json")
			b.UserInfo(rw, req)
			r, _ := io.ReadAll(rw.Body)
			response := &UserResponse{}
			json.Unmarshal(r, response)
			assert.Equal(t, tc.expectedAccountId, response.AccountId)
			assert.Equal(t, tc.expectedUserName, response.UserName)
		})
	}
}
