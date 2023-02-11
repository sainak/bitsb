package repo

import (
	"fmt"
	"strings"
)

type Filters map[string]interface{}

func (f *Filters) BuildQuery() string {
	var query []string
	for k, v := range *f {
		// keys are in format of "key:operator"
		s := strings.Split(k, ":")
		switch strings.ToLower(s[1]) {
		case "ilike", "like":
			v = fmt.Sprintf("%%%v%%", v)
		}
		query = append(query, fmt.Sprintf("%s %s '%v'", s[0], s[1], v))
	}
	return strings.Join(query, " AND ")
}
