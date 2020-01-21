package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Configuration struct {
	CookieSecret             string `yaml:"cookie_secret"`
	CsrfTokenValidTime       int    `yaml:"csrf_token_valid_time"`
	CsrfTokenSecret          string `yaml:"csrf_token_secret"`
	CsrfCookieName           string `yaml:"csrf_cookie_name"`
	RedisAddress             string `yaml:"redis_address"`
	UserInfoSessionValidTime int    `yaml:"userinfo_session_valid_time"`
	UserInfoSessionSecret    string `yaml:"userinfo_session_secret"`
	UserInfoSessionKey       string `yaml:"userinfo_session_key"`
	UserInfoCookieKey        string `yaml:"userinfo_cookie_key"`
}

const (
	ArticlesPerPage = 5
)

var configuration *Configuration

func LoadConfiguration(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	configuration = &config
	return err
}

func GetConfiguration() *Configuration {
	return configuration
}
