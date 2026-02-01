package configs

import "time"

type Auth struct {
	Domain string `yaml:"domain"`
	Token  Token  `yaml:"token"`
}

type Token struct {
	Secret              string        `yaml:"secret"`
	AccessTokenTTL      time.Duration `yaml:"accessTokenTTL"`
	RefreshTokenLongTTL time.Duration `yaml:"refreshTokenLongTTL"`
	RefreshTokenTTL     time.Duration `yaml:"refreshTokenTTL"`
}

func GetConfigAuth(c Config) Auth {
	return c.Auth
}

func GetConfigAuthToken(c Config) Token {
	return c.Auth.Token
}
