package plan

import _ "github.com/lib/pq"

type Postgres struct {
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	Address  string `hcl:"address"`
	Port     int    `hcl:"port"`
	Database string `hcl:"database"`
	Options  string `hcl:"options,optional"`
}
