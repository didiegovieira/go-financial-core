package account

type DeleteByPhoneInput struct {
	Phone string `form:"phone" binding:"required"`
}

type DeleteByPhoneOutput struct {
	Success bool `json:"success"`
}
