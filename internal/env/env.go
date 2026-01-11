package env

type EnvConfigs struct {
	EnvConfigs []EnvConfig `yaml:"envConfigs"`
}

type EnvConfig struct {
	Name   string `yaml:"name"`
	Folder string `yaml:"folder"`
}
