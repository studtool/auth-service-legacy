package queues

//go:generate easyjson

//easyjson:json
type RegistrationEmailData struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

//easyjson:json
type EmailUpdateData struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

//easyjson:json
type PasswordUpdateData struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
