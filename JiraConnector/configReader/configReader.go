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
	configReader.viperReader.SetConfigName("server")
	configReader.viperReader.SetConfigType("yaml")
	configReader.viperReader.AddConfigPath("../../ApiServer/configs")
	if err := configReader.viperReader.ReadInConfig(); err != nil {
		log.Fatal()
	}

	return &configReader
}

func (configReader *ConfigReader) GetLocalServerPort() uint {
	return configReader.viperReader.GetUint("connector_port")
}

func (configReader *ConfigReader) GetLocalServerHost() string {
	return configReader.viperReader.GetString("connector_host")
}

func (configReader *ConfigReader) GetThreadCount() int {
	return configReader.viperReader.GetInt("threadCount")
}

func (configReader *ConfigReader) GetJiraRepositoryUrl() string {
	return configReader.viperReader.GetString("jiraUrl")
}

func (configReader *ConfigReader) GetIssuesPerRequest() int {
	return configReader.viperReader.GetInt("issueInOneRequest")
}

func (configReader *ConfigReader) GetMinTimeSleep() int {
	return configReader.viperReader.GetInt("minTimeSleep")
}

func (configReader *ConfigReader) GetMaxTimeSleep() int {
	return configReader.viperReader.GetInt("maxTimeSleep")
}

func (configReader *ConfigReader) GetDbUsername() string {
	return configReader.viperReader.GetString("db_user")
}

func (configReader *ConfigReader) GetDbPassword() string {
	return configReader.viperReader.GetString("db_passwd")
}

func (configReader *ConfigReader) GetDbHost() string {
	return configReader.viperReader.GetString("db_host")
}

func (configReader *ConfigReader) GetDbPort() int {
	return configReader.viperReader.GetInt("db_port")
}

func (configReader *ConfigReader) GetDbName() string {
	return configReader.viperReader.GetString("db_name")
}
