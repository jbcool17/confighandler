package handler

type testConfig struct {
	FieldOne  int      `yaml:"fieldOne,omitempty"`
	Name      string   `yaml:"name"`
	Enabled   bool     `yaml:"enabled"`
	ListField []string `yaml:"listField,omitempty"`
	Options   options  `yaml:"options,omitempty"`
}

type testConfigs struct {
	TestConfigs []testConfig `yaml:"testConfigs"`
}

type options struct {
	Verbose bool
	Debug   bool
	Timeout int
	Version string
}
