package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Card struct {
	Name           string    `json:"name"`
	Number         string    `json:"number"`
	Balance        int       `json:"balance"`
	BalanceLocked  int       `json:"balanceLocked"`
	UpdateDateTime time.Time `json:"updateDateTime"`
}

type rawCards []Card

type prettyCard struct {
	Name           string
	Number         string
	Balance        string
	BalanceLocked  string
	UpdateDateTime string
}

func (c Card) PrettyPrint() (string, error) {
	pretty := prettyCard{
		Name:           c.Name,
		Number:         c.Number,
		Balance:        fmt.Sprintf("%.2fzł", float32(c.Balance)/100),
		BalanceLocked:  fmt.Sprintf("%.2fzł", float32(c.BalanceLocked)/100),
		UpdateDateTime: c.UpdateDateTime.Format("2006-01-02 15:04:05"),
	}

	if out, err := json.MarshalIndent(pretty, "", "  "); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}

func GetCardDetails(token string) (Card, error) {
	url := "https://api4you.sodexo.pl/api/card"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Card{}, errors.Wrap(err, "while creating request")
	}

	req.Header.Add("authority", "api4you.sodexo.pl")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("dnt", "1")
	req.Header.Add("authorization", token)
	req.Header.Add("accept-language", "pl")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://dlaciebie.sodexo.pl")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://dlaciebie.sodexo.pl/")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Card{}, errors.Wrap(err, "while sending request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Card{}, errors.Wrap(err, "while reading response body")
	}

	var cardsResp rawCards
	if err := json.Unmarshal(body, &cardsResp); err != nil {
		return Card{}, errors.Wrap(err, "while unmarshalling body to structure")
	}

	if len(cardsResp) == 0 {
		return Card{}, errors.New("no Sodexo Cards sent by api")
	}

	if len(cardsResp) > 1 {
		log.Println("You have more than 1 Sodexo Card, handling the first one")
	}

	return cardsResp[0], nil
}
