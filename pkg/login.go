package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type loginResponse struct {
	Token string `json:"token"`
}

func Login(login, password string) (string, error) {
	url := "https://api4you.sodexo.pl/api/user/login"

	payload := strings.NewReader(fmt.Sprintf(`{"login":"%s","password":"%s","deviceData":{"deviceOrigin":"WEB"}}`, login, password))

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return "", errors.New("while creating request")
	}

	req.Header.Add("authority", "api4you.sodexo.pl")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("dnt", "1")
	req.Header.Add("authorization", "")
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
		return "", errors.Wrap(err, "while sending request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "while reading response body")
	}

	var loginResp loginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		defer fmt.Printf("Logging response body for debugging: %s\n", string(body))
		return "", errors.Wrap(err, "while unmarshalling body to structure")
	}

	return loginResp.Token, nil
}
