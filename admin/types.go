package admin

import "github.com/jasonsnider/com.jasonsnider.go/internal/types"

type UserRegistrationTemplate struct {
	Title            string
	Description      string
	Keywords         string
	Body             string
	ValidationErrors map[string]string
	User             types.RegisterUser
	BustCssCache     string
	BustJsCache      string
}
