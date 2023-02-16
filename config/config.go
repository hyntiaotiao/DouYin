package config

var (
	MYSQL_CONFIG  = &MysqlConfig{Timeout: "10s", MaxOpenConns: 100, MaxIdleConns: 20}
	SERVER_CONFIG = &ServerConfig{}
)

type MysqlConfig struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Dbname       string `json:"dbName"`
	Timeout      string `json:"timeout"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

type ServerConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}
