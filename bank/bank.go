package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Bank interface {
	UserInfo(rw http.ResponseWriter, req *http.Request)
	RegisterHandlers()
}

type bank struct {
	registeredUsers map[string]*user
}

func NewBank() Bank {
	um := NewUserManager()
	oI := OnboartingInfo()
	ru := um.OnboardUsers(oI)
	return &bank{
		registeredUsers: ru,
	}
}

func (b *bank) RegisterHandlers() {
	http.HandleFunc("/userInfo", b.UserInfo)
}

func (b *bank) getUser(userName string) (*user, error) {
	//var currentUser *user
	currentUser, ok := b.registeredUsers[userName]
	if !ok {
		return nil, fmt.Errorf("user name %v is not a valid user", userName)
	}
	return currentUser, nil
}

func (b *bank) UserInfo(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method != "GET" {
		rw.Write([]byte("invalid method"))
		rw.WriteHeader(429)
		return
	}
	qV := req.URL.Query()
	if len(qV["userName"]) == 0 {
		rw.Write([]byte("query parameter missing"))
		rw.WriteHeader(400)
		return
	}
	userName := qV["userName"][0]
	if userName == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("missing query parameter"))
		return
	}
	u, err := b.getUser(userName)
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte(err.Error()))
		return
	}
	respBody := &UserResponse{
		UserName:       u.userName,
		AccountId:      u.accountInfo.accountId,
		CurrentBalance: u.accountInfo.currentBalance,
	}
	body, _ := json.Marshal(respBody)
	rw.Write(body)
}
