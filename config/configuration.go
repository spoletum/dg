package config

type Configuration struct {
	Postgres []PostgresBlock `hcl:"postgres,block"`
}
