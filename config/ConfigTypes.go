package config
type RedisConfig struct {
	Host     string
	Password string
	Db       int
}
type MongoConfig struct {
	Connection string
	DbName string
}


