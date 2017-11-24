package dialect

import (
	"fmt"

	"github.com/posener/orm/def"
	"github.com/posener/orm/dialect/api"
	"github.com/posener/orm/dialect/sqlite3"
)

func New(name string, tp def.Type) (api.Dialect, error) {
	switch name {
	case "sqlite3", "sqlite":
		return sqlite3.New(tp), nil
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", name)
	}
}
