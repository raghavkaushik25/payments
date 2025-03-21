package bank

type UserResponse struct {
	UserName       string `json:"user_name"`
	AccountId      string `json:"account_id"`
	CurrentBalance int    `json:"current_balance"`
}

type TransferRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type TransferResponse struct {
	Message         string `json:"message"`
	UpdatedBalance  int    `json:"updated_balance"`
	PreviousBalance int    `json:"previous_balance"`
}
