package item

type Type string
type Item interface {
	Path() string
	Type() Type
}

const (
	FILE Type = "File"
	DIR  Type = "Directory"
)
