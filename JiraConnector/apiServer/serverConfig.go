package apiServer

type ServerConfig struct {
	port uint
	host string
}

func NewServerConfig(host string, port uint) *ServerConfig {
	return &ServerConfig{
		port: port,
		host: host,
	}
}
