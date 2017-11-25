package dialect

import (
	"fmt"

	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/sqlite3"
)

// New returns a new Dialect if one is implemented
func New(name string, tp common.Type) (common.Dialect, error) {
	switch name {
	case "sqlite3", "sqlite":
		return sqlite3.New(tp), nil
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", name)
	}
}
