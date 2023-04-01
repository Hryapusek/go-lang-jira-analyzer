package configReader

import (
	"github.com/spf13/viper"
	"log"
)

type ConfigReader struct {
	viperReader *viper.Viper
}

func NewConfigReader() *ConfigReader {
	configReader := ConfigReader{}
	configReader.viperReader = viper.New()
	configReader.viperReader.SetConfigName("config")
	configReader.viperReader.SetConfigType("yaml")
	configReader.viperReader.AddConfigPath("./configurationFiles")
	if err := configReader.viperReader.ReadInConfig(); err != nil {
		log.Fatal()
	}

	return &configReader
}

func (configReader *ConfigReader) GetLocalServerPort() uint {
	return configReader.viperReader.GetUint("ProgramSettings.local_http_server_port")
}

func (configReader *ConfigReader) GetThreadCount() uint {
	return configReader.viperReader.GetUint("ProgramSettings.threadCount")
}

func (configReader *ConfigReader) GetJiraRepositoryUrl() string {
	return configReader.viperReader.GetString("ProgramSettings.jiraUrl")
}

func (configReader *ConfigReader) GetIssuesPerRequest() uint {
	return configReader.viperReader.GetUint("ProgramSettings.issueInOneRequest")
}
