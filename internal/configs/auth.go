package configs

type Auth struct {
	AccessToken AccessToken `yaml:"accessToken"`
}

type AccessToken struct {
	ExpiresIn string `yaml:"expiresIn"`
}
