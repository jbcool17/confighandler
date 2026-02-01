package pkg

type Cluster struct {
	Name          string        `yaml:"name"`
	Location      string        `yaml:"location"`
	NetworkConfig NetworkConfig `yaml:"network_config"`
	NodePools     []NodePool    `yaml:"node_pools"`
	FeatureOne    Enabled       `yaml:"feature_one"`
	FeatureTwo    Enabled       `yaml:"feature_two"`
	Test          string        `yaml:"test"`
}

type Enabled struct {
	Enabled bool `yaml:"enabled"`
}

type NetworkConfig struct {
	Network string `yaml:"network"`
	PodCIDR string `yaml:"pod_cidr"`
}

type NodePool struct {
	Name        string `yaml:"name"`
	MachineType string `yaml:"machine_type"`
}
