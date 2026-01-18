package configs

type CacheType string

const (
	CacheTypeRam   CacheType = "ram"
	CacheTypeRedis CacheType = "redis"
)

type Cache struct {
	Type     CacheType `yaml:"type"`
	Address  string    `yaml:"address"`
	Port     string    `yaml:"port"`
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
}
