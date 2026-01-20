package configs

type Grpc struct {
	AccountService AccountService `yaml:"account_service"`
}

type AccountService struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

func GetConfigAccountService(c Config) AccountService {
	return c.Grpc.AccountService
}
