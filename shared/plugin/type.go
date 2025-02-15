package plugin

import "strings"

type Type int

const (
	Undefined Type = iota
	Logger
)

func (t Type) String() string {
	return [...]string{"Undefined", "Logger"}[t]
}

func Parse(val string) Type {
	switch strings.ToLower(val) {
	case "logger":
		return Logger
	default:
		return Undefined
	}

}
