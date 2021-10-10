package site

import (
	"encoding/json"
	"fmt"
	"github.com/loginradius/lr-cli/api"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

type TokenResp struct {
	AppName string `json:"app_name"`
	XSign   string `json:"xsign"`
	XToken  string `json:"xtoken"`
}

func TestSetSite(t *testing.T) {

	prevAppId := appid

	user, _ := user.Current()

	baseFileName := filepath.Join(user.HomeDir, ".lrcli")

	defer func() {
		appid = prevAppId
	}()

	tests := []struct {
		name    string
		args    map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			"invalid app id",
			map[string]interface{}{
				"appId":          int64(123456),
				"pathToSiteInfo": "siteInfo.json",
				"pathToToken":    "token.json",
				"function":       NewSiteCmd,
			},
			"There is no site with this AppID.\n",
			false,
		},
		{
			"invalid app id",
			map[string]interface{}{
				"appId":          int64(0),
				"pathToSiteInfo": "siteInfo.json",
				"pathToToken":    "token.json",
				"function":       setSite,
			},
			"You are already using this site\n",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			appid = tt.args["appId"].(int64)

			createSiteInfo(tt, baseFileName)

			createToken(tt, baseFileName)

			output := captureOutput(tt.args["function"])
			assert.Equal(t, tt.want, output)

			removeFile(baseFileName)

		})
	}
}

func captureOutput(fun interface{}) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	switch f := fun.(type) {
	case func() *cobra.Command:
		f().Execute()
	case func() error:
		f()
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func createToken(tt struct {
	name    string
	args    map[string]interface{}
	want    string
	wantErr bool
}, baseFileName string) error {
	token := TokenResp{
		"app_name",
		"sign",
		"token",
	}

	fileName := filepath.Join(baseFileName, tt.args["pathToToken"].(string))

	dest, err := os.Create(fileName)
	if err != nil {
		return err
	}

	bytes, _ := json.Marshal(token)

	fmt.Fprintf(dest, string(bytes))

	return nil
}

func createSiteInfo(tt struct {
	name    string
	args    map[string]interface{}
	want    string
	wantErr bool
}, baseFileName string) error {

	err := os.Mkdir(baseFileName, 0755)
	if err != nil {
		fmt.Println("os.Create:", err)
	}

	fileName := filepath.Join(baseFileName, tt.args["pathToSiteInfo"].(string))

	dest, err := os.Create(fileName)

	if err != nil {
		return err
	}

	siteInfo := map[int64]api.SitesReponse{
		tt.args["appId"].(int64): {
			Appname: "app_name",
		},
	}

	bytes, _ := json.Marshal(siteInfo)

	fmt.Fprintf(dest, string(bytes))

	return nil
}

func removeFile(baseFileName string) {
	err := os.RemoveAll(baseFileName)
	if err != nil {
		fmt.Println("os.Create:", err)
	}
}