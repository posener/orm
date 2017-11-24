package sqlite3

import (
	"regexp"
	"strings"
	"testing"

	"github.com/posener/orm/dialect/api"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		columner columner
		wherer   wherer
		groups   []api.Group
		orders   []api.Order
		pager    pager
		want     string
	}{
		{
			want: "SELECT * FROM 'name'",
		},
		{
			columner: columner{},
			want:     "SELECT * FROM 'name'",
		},
		{
			columner: columner{},
			pager:    pager{},
			want:     "SELECT * FROM 'name'",
		},
		{
			columner: columner{count: true},
			want:     "SELECT COUNT(*) FROM 'name'",
		},
		{
			columner: columner{columns: []string{"a", "b", "c"}, count: true},
			want:     "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name'",
		},
		{
			columner: columner{columns: []string{"a", "b", "c"}},
			want:     "SELECT `a`, `b`, `c` FROM 'name'",
		},
		{
			pager: pager{},
			want:  "SELECT * FROM 'name'",
		},
		{
			pager: pager{limit: 1},
			want:  "SELECT * FROM 'name' LIMIT 1",
		},
		{
			pager: pager{limit: 1, offset: 2},
			want:  "SELECT * FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			pager: pager{offset: 1},
			want:  "SELECT * FROM 'name' LIMIT 0 OFFSET 1",
		},
		{
			columner: columner{columns: []string{"a", "b", "c"}, count: true},
			pager:    pager{limit: 1, offset: 2},
			want:     "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name' LIMIT 1 OFFSET 2",
		},
		{
			groups: []api.Group{{Column: "a"}, {Column: "b"}},
			want:   "SELECT * FROM 'name' GROUP BY `a`, `b`",
		},
		{
			orders: []api.Order{
				{Column: "c", Dir: "ASC"},
				{Column: "d", Dir: "DESC"},
			},
			want: "SELECT * FROM 'name' ORDER BY `c` ASC, `d` DESC",
		},
		{
			wherer: wherer("`k` > 0"),
			want:   "SELECT * FROM 'name' WHERE `k` > 0",
		},
		{
			columner: columner{columns: []string{"a", "b", "c"}, count: true},
			wherer:   wherer("`k` > 0"),
			groups:   []api.Group{{Column: "a"}, {Column: "b"}},
			orders: []api.Order{
				{Column: "c", Dir: "ASC"},
				{Column: "d", Dir: "DESC"},
			},
			pager: pager{limit: 1, offset: 2},
			want:  "SELECT `a`, `b`, `c`, COUNT(*) FROM 'name' WHERE `k` > 0 GROUP BY `a`, `b` ORDER BY `c` ASC, `d` DESC LIMIT 1 OFFSET 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Select(tabler{}, &tt.columner, &tt.wherer, tt.groups, tt.orders, &tt.pager)
			assert.Equal(t, tt.want, reduceSpaces(got), " ")
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign []api.Assignment
		want   string
	}{
		{
			want: "INSERT INTO 'name' () VALUES ()",
		},
		{
			assign: []api.Assignment{{Column: "c1", Value: 1}},
			want:   "INSERT INTO 'name' (`c1`) VALUES (?)",
		},
		{
			assign: []api.Assignment{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			want:   "INSERT INTO 'name' (`c1`, `c2`) VALUES (?, ?)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Insert(tabler{}, tt.assign)
			assert.Equal(t, tt.want, reduceSpaces(got), " ")
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		assign []api.Assignment
		wherer wherer
		want   string
	}{
		{
			want: "UPDATE 'name' SET",
		},
		{
			assign: []api.Assignment{{Column: "c1", Value: 1}},
			want:   "UPDATE 'name' SET `c1` = ?",
		},
		{
			assign: []api.Assignment{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			want:   "UPDATE 'name' SET `c1` = ?, `c2` = ?",
		},
		{
			assign: []api.Assignment{{Column: "c1", Value: 1}, {Column: "c2", Value: ""}},
			wherer: wherer("`k` > 3"),
			want:   "UPDATE 'name' SET `c1` = ?, `c2` = ? WHERE `k` > 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Update(tabler{}, tt.assign, &tt.wherer)
			assert.Equal(t, tt.want, reduceSpaces(got), " ")
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wherer wherer
		want   string
	}{
		{
			want: "DELETE FROM 'name'",
		},
		{
			wherer: wherer("`k` > 3"),
			want:   "DELETE FROM 'name' WHERE `k` > 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Delete(tabler{}, &tt.wherer)
			assert.Equal(t, tt.want, reduceSpaces(got), " ")
		})
	}
}

func reduceSpaces(s string) string {
	re := regexp.MustCompile("([ ]+)")
	return strings.Trim(re.ReplaceAllString(s, " "), " ")
}

type columner struct {
	columns []string
	count   bool
}

func (c *columner) Columns() []string {
	return c.columns
}

func (c *columner) Count() bool {
	return c.count
}

type wherer string

func (w *wherer) Where() string {
	return string(*w)
}

type pager struct {
	limit, offset int64
}

func (p *pager) Page() (int64, int64) {
	return p.limit, p.offset
}

type tabler struct{}

func (tabler) Table() string {
	return "name"
}
