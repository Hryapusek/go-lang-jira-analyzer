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

func (configReader *ConfigReader) GetThreadCount() int {
	return configReader.viperReader.GetInt("ProgramSettings.threadCount")
}

func (configReader *ConfigReader) GetJiraRepositoryUrl() string {
	return configReader.viperReader.GetString("ProgramSettings.jiraUrl")
}

func (configReader *ConfigReader) GetIssuesPerRequest() int {
	return configReader.viperReader.GetInt("ProgramSettings.issueInOneRequest")
}

func (configReader *ConfigReader) GetMinTimeSleep() int {
	return configReader.viperReader.GetInt("ProgramSettings.minTimeSleep")
}

func (configReader *ConfigReader) GetMaxTimeSleep() int {
	return configReader.viperReader.GetInt("ProgramSettings.maxTimeSleep")
}

func (configReader *ConfigReader) GetDbUsername() string {
	return configReader.viperReader.GetString("DBSettings.username")
}

func (configReader *ConfigReader) GetDbPassword() string {
	return configReader.viperReader.GetString("DBSettings.password")
}

func (configReader *ConfigReader) GetDbHost() string {
	return configReader.viperReader.GetString("DBSettings.hostname")
}

func (configReader *ConfigReader) GetDbPort() int {
	return configReader.viperReader.GetInt("DBSettings.port")
}

func (configReader *ConfigReader) GetDbName() string {
	return configReader.viperReader.GetString("DBSettings.name")
}
