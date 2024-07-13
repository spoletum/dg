package plan

// Storage represents where the snapshot will be stored
type Storage struct {
	Type string `hcl:"type,label"`
	Path string `hcl:"path"`
}
