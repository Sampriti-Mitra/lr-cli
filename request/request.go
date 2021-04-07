package request

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
)

func Rest(method string, url string, headers map[string]string, payload string) ([]byte, error) {
	conf := config.GetInstance()
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		log.Printf("error while Performing the Request: %s", err.Error())
		return nil, err
	}

	type TokenResp struct {
		AppName string `json:"app_name"`
		XSign   string `json:"xsign"`
		XToken  string `json:"xtoken"`
	}

	var token TokenResp

	// LoginRadius Default Headers
	v2, err := cmdutil.GetCreds()
	err = json.Unmarshal(v2, &token)
	if err == nil && token.AppName != "" {
		req.Header.Set("x-is-loginradius--sign", token.XSign)
		req.Header.Set("x-is-loginradius--token", token.XToken)
	} else if !strings.Contains(url, "auth/login") {
		return nil, errors.New("Please Login to execute this command")
	}
	req.Header.Set("Origin", conf.DashboardDomain)
	req.Header.Set("x-is-loginradius-ajax", "true")

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("%s", err.Error())
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
