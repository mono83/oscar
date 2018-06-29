package impact

import (
	"errors"
	"strings"
)

// Parse parses incoming string as impact level
func Parse(in string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "0", "none", "no":
		return None, nil
	case "1", "read":
		return Read, nil
	case "3", "create":
		return Create, nil
	case "4", "write", "modify", "critical":
		return Modify, nil
	}

	return None, errors.New("unknown impact level " + in)
}

// ParseOrDefault parses incoming string as impact level and returns Default on error
func ParseOrDefault(in string) Level {
	l, err := Parse(in)
	if err != nil {
		l = Default
	}

	return l
}
