package cmd

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/vrischmann/envconfig"

	"github.com/aerfio/sdxo/pkg"
)

type config struct {
	Login    string
	Password string
}

func Run() error {
	cfg := config{}
	err := envconfig.InitWithPrefix(&cfg, "SDXO")
	if err != nil {
		return errors.Wrap(err, "while reading login and password from env vars")
	}
	token, err := pkg.Login(cfg.Login, cfg.Password)
	if err != nil {
		return errors.Wrap(err, "while logging in")
	}

	card, err := pkg.GetCardDetails(token)
	if err != nil {
		return errors.Wrap(err, "while requesting card details")
	}

	cardString, err := card.PrettyPrint()
	if err != nil {
		return errors.Wrap(err, "while pretty printing card details")
	}

	fmt.Println(cardString)

	return nil
}
