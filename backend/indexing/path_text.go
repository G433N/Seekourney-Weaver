package indexing

import (
	"seekourney/core/database"
	"seekourney/utils"
)

type PathText struct {
	Path utils.Path
	Text string
}

func (pathText PathText) SQLGetName() string {
	return "path_text"
}

func (pathText PathText) SQLGetFields() []string {
	return []string{"path", "plain_text"}
}

func (pathText PathText) SQLGetValues() []any {
	return []database.SQLValue{pathText.Path, pathText.Text}
}
