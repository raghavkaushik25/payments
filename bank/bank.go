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
	http.HandleFunc("/transfer", b.Transfer)
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

func (b *bank) Transfer(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.WriteHeader(429)
		rw.Write([]byte("invalid method"))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	transfer := &TransferRequest{}
	transferBody := &TransferResponse{}
	err := json.NewDecoder(req.Body).Decode(transfer)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("internal server error %v", err)))
		return
	}

	from, err := b.getUser(transfer.From)
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte(err.Error()))
		return
	}
	from.accountInfo.l.Lock()
	defer from.accountInfo.l.Unlock()
	to, err := b.getUser(transfer.To)
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte(err.Error()))
		return
	}
	to.accountInfo.l.Lock()
	defer to.accountInfo.l.Unlock()
	if transfer.Amount > from.accountInfo.currentBalance {
		rw.WriteHeader(400)
		rw.Write([]byte(fmt.Sprintf("insufficient balance %v", from.accountInfo.currentBalance)))
		return
	}
	if from.accountInfo.accountId == to.accountInfo.accountId {
		rw.WriteHeader(400)
		rw.Write([]byte(fmt.Sprintf("cannot transfer to the same account %v", from.accountInfo.accountId)))
		return
	}
	transferBody.PreviousBalance = from.accountInfo.currentBalance
	from.accountInfo.currentBalance -= transfer.Amount
	to.accountInfo.currentBalance += +transfer.Amount
	transferBody.UpdatedBalance = from.accountInfo.currentBalance
	transferBody.Message = fmt.Sprintf("accound Id : %v has been debited with ammount %v; updated balance is %v", from.accountInfo.accountId, transfer.Amount, from.accountInfo.currentBalance)
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(transferBody)
}
