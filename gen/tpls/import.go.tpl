package {{$.Type.Package}}

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/posener/orm"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect"
	{{ range $_, $import := .Type.Imports -}}
    "{{$import}}"
    {{ end -}}
)

