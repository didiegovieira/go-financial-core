package account

type FindByPhoneInput struct {
	Phone string `form:"phone" binding:"required"`
}

type FindByPhoneOutput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	DocumentType    string `json:"document_type"`
	DocumentNumber  string `json:"document_number"`
	DocumentCountry string `json:"document_country"`
}
