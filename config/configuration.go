package config

type Configuration struct {
	Postgres []Postgres `hcl:"postgres,block"`
}
