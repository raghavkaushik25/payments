package bank

import (
	"bytes"
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

func TestTransfer(t *testing.T) {
	um := NewUserManager()
	oI := OnboartingInfo()
	ru := um.OnboardUsers(oI)
	b := bank{
		registeredUsers: ru,
	}
	type tests struct {
		name     string
		tr       *TransferRequest
		expected interface{}
	}
	tcs := []tests{
		{
			name: "fromMarkToAdam",
			tr: &TransferRequest{
				From:   "Mark",
				To:     "Adam",
				Amount: 1,
			},
			expected: 99,
		}, {
			name: "fromMarkToAdam",
			tr: &TransferRequest{
				From:   "Mark",
				To:     "Adam",
				Amount: 5,
			},
			expected: 94,
		}, {
			name: "fromMarkToAdam",
			tr: &TransferRequest{
				From:   "Mark",
				To:     "Adam",
				Amount: 4,
			},
			expected: 90,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			by, _ := json.Marshal(tc.tr)
			reader := bytes.NewReader(by)
			req := httptest.NewRequest("POST", "http://localhost:8080/transfer", reader)
			req.Header.Set("Content-Type", "application/json")
			b.Transfer(rw, req)
			r, _ := io.ReadAll(rw.Body)
			response := &TransferResponse{}
			json.Unmarshal(r, response)
			assert.Equal(t, tc.expected, response.UpdatedBalance)
		})
	}
}

// Tests transactions going to the same account concurrently.
func TestConcurrentTransfer(t *testing.T) {
	um := NewUserManager()
	oI := OnboartingInfo()
	ru := um.OnboardUsers(oI)
	b := bank{
		registeredUsers: ru,
	}
	type tests struct {
		name string
		tr   *TransferRequest
	}
	tcs := []tests{
		{
			name: "fromMarkToAdam#1",
			tr: &TransferRequest{
				From:   "Mark",
				To:     "Adam",
				Amount: 10,
			},
		},
		{
			name: "fromJaneToAdam#1",
			tr: &TransferRequest{
				From:   "Jane",
				To:     "Adam",
				Amount: 10,
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rw := httptest.NewRecorder()
			by, _ := json.Marshal(tc.tr)
			reader := bytes.NewReader(by)
			req := httptest.NewRequest("POST", "http://localhost:8080/transfer", reader)
			req.Header.Set("Content-Type", "application/json")
			b.Transfer(rw, req)
			r, _ := io.ReadAll(rw.Body)
			response := &TransferResponse{}
			json.Unmarshal(r, response)
		})
	}
	t.Cleanup(func() {
		newBalanceMark, ok := b.registeredUsers["Mark"]
		assert.Equal(t, ok, true)
		assert.Equal(t, 90, newBalanceMark.accountInfo.currentBalance)
		newBalanceAdam, ok := b.registeredUsers["Adam"]
		assert.Equal(t, ok, true)
		assert.Equal(t, 20, newBalanceAdam.accountInfo.currentBalance)
		newBalanceJane, ok := b.registeredUsers["Jane"]
		assert.Equal(t, ok, true)
		assert.Equal(t, 40, newBalanceJane.accountInfo.currentBalance)
	})
}
