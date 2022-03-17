package listener

import (
	"encoding/json"
	"fmt"
	"time"

	dq "doublequote/pkg/config"
	"doublequote/pkg/domain"
)

func sendVerificationEmail(e domain.Event, emailService domain.EmailService, cryptoService domain.CryptoService, cfg dq.Config) error {
	var p domain.UserCreatedPayload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return err
	}

	token, err := cryptoService.CreateToken(map[string]interface{}{"id": p.User.ID}, (time.Hour*24)*2)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/verify?t=%s", cfg.FrontendURL(), token)

	err = emailService.SendEmail(
		[]string{p.User.Email},
		"Verify your Doublequote email",
		"Verify your Doublequote email by clicking this link: "+url,
	)
	return err
}
