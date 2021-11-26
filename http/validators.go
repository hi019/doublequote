package http

import (
	"context"
	"log"

	dq "doublequote"
	"github.com/gobuffalo/validate"
)

type EmailIsUnique struct {
	Email       string
	UserService dq.UserService
}

func (v *EmailIsUnique) IsValid(errors *validate.Errors) {
	_, c, err := v.UserService.FindUsers(context.Background(), dq.UserFilter{Email: &v.Email, Limit: 1}, dq.UserInclude{})

	if err != nil {
		log.Println("[emailIsUnique] error: " + err.Error())
	}

	if c > 0 {
		errors.Add("email", dq.ErrEmailTaken)
	}
}
