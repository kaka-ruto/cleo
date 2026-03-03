package qacatalog

type Actor struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Surfaces    []string `yaml:"surfaces"`
	AuthProfile string   `yaml:"auth_profile"`
}
