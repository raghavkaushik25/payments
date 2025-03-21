package bank

import (
	"github.com/google/uuid"
)

type Unlock func()

type account struct {
	accountId      string
	currentBalance int
	openingBalance int
}

type user struct {
	userName    string
	userId      string
	accountInfo *account
}

func NewUserManager() *user {
	return &user{}
}

func (u *user) OnboardUsers(oI map[string]interface{}) map[string]*user {
	RegisteredUsers := map[string]*user{}
	for k, v := range oI {
		user := &user{
			userName:    k,
			userId:      uuid.NewString(),
			accountInfo: u.onboardUserAccount(v.(map[string]interface{})["initialBalance"].(int)),
		}
		RegisteredUsers[user.userName] = user
	}
	return RegisteredUsers
}

func (u *user) onboardUserAccount(openingBalance int) *account {
	return &account{
		accountId:      uuid.NewString(),
		openingBalance: openingBalance,
		currentBalance: openingBalance,
	}
}

func OnboartingInfo() map[string]interface{} {
	oI := make(map[string]interface{})
	oI["Mark"] = map[string]interface{}{
		"initialBalance": 100,
	}
	oI["Jane"] = map[string]interface{}{
		"initialBalance": 50,
	}
	oI["Adam"] = map[string]interface{}{
		"initialBalance": 0,
	}
	return oI
}
