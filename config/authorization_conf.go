package config

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type AuthorizationConfig struct {
	user     string
	password string
}

func (this *AuthorizationConfig) SetupUser(user, password string) {
	this.user = user
	this.password = password
}

func (this *AuthorizationConfig) ValidationForBasicAuthorization(value string) bool {
	if len(value) < 7 {
		return false
	}
	if !strings.HasPrefix(value, "Basic ") {
		return false
	}
	value = strings.ReplaceAll(value, "Basic ", "")
	strByte, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pair := strings.Split(string(strByte), ":")
	if len(pair) != 2 {
		return false
	}
	if this.user == pair[0] && this.password == pair[1] {
		return true
	}
	return false
}
