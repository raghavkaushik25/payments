package bank

type UserResponse struct {
	UserName       string `json:"user_name"`
	AccountId      string `json:"account_id"`
	CurrentBalance int    `json:"current_balance"`
}
