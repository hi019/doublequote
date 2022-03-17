package http

import (
	"context"
	"log"

	"doublequote/pkg/domain"
	"github.com/gobuffalo/validate"
)

type EmailIsUnique struct {
	Email       string
	UserService domain.UserService
}

func (v *EmailIsUnique) IsValid(errors *validate.Errors) {
	_, c, err := v.UserService.FindUsers(context.Background(), domain.UserFilter{Email: &v.Email, Limit: 1}, domain.UserInclude{})

	if err != nil {
		log.Println("[emailIsUnique] error: " + err.Error())
	}

	if c > 0 {
		errors.Add("email", domain.ErrEmailTaken)
	}
}
