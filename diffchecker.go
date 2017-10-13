package diffchecker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	DIFFCHECKERURL        = "https://diffchecker-api-production.herokuapp.com"
	DIFFCHECKERSESSIONS   = DIFFCHECKERURL + "/sessions"
	DIFFCHECKERDIFFS      = DIFFCHECKERURL + "/diffs"
	DIFFCHECKERSUCCESSURL = "https://www.diffchecker.com/"
	AUTHTOKENKEY          = "authToken"
)

type DiffChecker struct {
	email    string
	password string
}

func (checker DiffChecker) Upload(left string, right string, title string) (slug string, err error) {
	return checker.UploadWithDuration(left, right, title, FOREVER)
}

func (checker DiffChecker) UploadBytes(left []byte, right []byte, title string) (slug string, err error) {
	return checker.UploadBytesWithDuration(left, right, title, FOREVER)
}

func (checker DiffChecker) UploadBytesWithDuration(left []byte, right []byte, title string, expiry DiffCheckerExpiry) (slug string, err error) {
	return checker.UploadWithDuration(string(left), string(right), title, expiry)
}

func (checker DiffChecker) UploadWithDuration(left string, right string, title string, expiry DiffCheckerExpiry) (slug string, err error) {
	token, err := checker.auth()

	if err != nil {
		return "", err
	}

	urlValues := url.Values{
		"left":   {left},
		"right":  {right},
		"expiry": {expiry.String()},
	}

	if len(title) > 0 {
		urlValues.Add("title", title)
	}

	request, _ := http.NewRequest("POST", DIFFCHECKERDIFFS, strings.NewReader(urlValues.Encode()))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)
	defer response.Body.Close()

	if response.StatusCode != 201 {
		return "", fmt.Errorf("response was %d, not 201 as expected", response.StatusCode)
	}

	jsonBody, err := parseJson(response.Body)

	if err != nil {
		return "", err
	}

	return DIFFCHECKERSUCCESSURL + jsonBody["slug"].(string), nil
}

func (checker DiffChecker) auth() (token string, err error) {
	response, err := http.PostForm(DIFFCHECKERSESSIONS, url.Values{"email": {checker.email}, "password": {checker.password}})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("response was %d, not 200 as expected", response.StatusCode)
	}

	jsonBody, err := parseJson(response.Body)

	if err != nil {
		return "", err
	}

	if jsonBody[AUTHTOKENKEY] == nil {
		return "", fmt.Errorf("response did not contain %s", AUTHTOKENKEY)
	}

	return jsonBody[AUTHTOKENKEY].(string), nil
}

func parseJson(body io.ReadCloser) (jsonBody map[string]interface{}, err error) {
	bytes, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &jsonBody)

	return
}