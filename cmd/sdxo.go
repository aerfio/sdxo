package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/aerfio/sdxo/pkg"
)

type config struct {
	Login    string
	Password string
}

func setupConfig() (config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/sdxo")
	viper.AddConfigPath("$HOME/.sdxo")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return config{}, errors.Wrap(err, "configuration file config.yaml not found")
		} else {
			return config{}, errors.Wrap(err, "reading configuration")
		}
	}

	var c config
	if err := viper.Unmarshal(&c); err != nil {
		return config{}, err
	}

	return c, nil
}

func Run() error {
	cfg, err := setupConfig()
	if err != nil {
		return errors.Wrap(err, "while reading login and password from config file")
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
