package config

type GatewayConfig struct {
	GatewayHost      string `yaml:"host"`
	GatewayPort      int    `yaml:"port"`
	GatewayAPIPrefix string `yaml:"api_prefix"`

	ResourceHost      string `yaml:"resource_host"`
	ResourcePort      int    `yaml:"resource_port"`
	ResourceAPIPrefix string `yaml:"resource_api_prefix"`
	ResourceTimeout   int    `yaml:"resourceTimeout"`

	AnalyticsHost      string `yaml:"analytics_host"`
	AnalyticsPort      int    `yaml:"analytics_port"`
	AnalyticsAPIPrefix string `yaml:"analytics_api_prefix"`
	AnalyticsTimeout   int    `yaml:"analyticsTimeout"`

	ConnectorHost      string `yaml:"connector_host"`
	ConnectorPort      int    `yaml:"connector_port"`
	ConnectorAPIPrefix string `yaml:"connector_api_prefix"`
}

type ResourceConfig struct {
	ResourceHost      string `yaml:"resource_host"`
	ResourcePort      int    `yaml:"resource_port"`
	ResourceAPIPrefix string `yaml:"resource_api_prefix"`
	MainAPIPrefix     string `yaml:"api_prefix"`
}

type AnalyticsConfig struct {
	AnalyticsHost      string `yaml:"analytics_host"`
	AnalyticsPort      int    `yaml:"analytics_port"`
	AnalyticsAPIPrefix string `yaml:"analytics_api_prefix"`
	MainAPIPrefix      string `yaml:"api_prefix"`
}

type ConnectorConfig struct {
	ConnectorHost      string `yaml:"connector_host"`
	ConnectorPort      int    `yaml:"connector_port"`
	ConnectorAPIPrefix string `yaml:"connector_api_prefix"`
	MainAPIPrefix      string `yaml:"api_prefix"`
}
