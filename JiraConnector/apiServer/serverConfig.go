package apiServer

type ServerConfig struct {
	port uint
}

func NewServerConfig(port uint) *ServerConfig {
	return &ServerConfig{
		port: port,
	}
}
