package service

type VacConfig struct {
	Name string
	Bind string
	Log  string
	Auth struct {
		Enabled  bool   `yaml: "enabled"`
		User     string `yaml: "user"`
		Password string `yaml: "password"`
	} `yaml: "auth"`
	InAllowList  []string `yaml:"inAllowList"`
	OutAllowList []string `yaml:"outAllowList"`
}
