package format

import (
	"testing"

	"strings"

	"github.com/posener/orm/common"
	"github.com/stretchr/testify/assert"
)

func TestColumnsJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		p        common.SelectParams
		wantCols string
		wantJoin string
	}{
		{
			p:        common.SelectParams{Table: "table", Columns: selector{}},
			wantCols: "`table`.*",
		},
		{
			p:        common.SelectParams{Table: "table", Columns: selector{count: true}},
			wantCols: "COUNT(*)",
		},
		{
			p:        common.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}}},
			wantCols: "`table`.`a`, `table`.`b`",
		},
		{
			p:        common.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}, count: true}},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p:        common.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}, count: true}},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{
								Table:   "B",
								Columns: selector{},
							},
							Pairings: []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.*, `B`.*",
			wantJoin: "JOIN (`B`) ON (`A`.`B_id` = `B`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					count: true,
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{Table: "B", Columns: selector{}},
							Pairings:     []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`B`.*, COUNT(*)",
			wantJoin: "JOIN (`B`) ON (`A`.`B_id` = `B`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{Table: "B", Columns: selector{}},
							Pairings:     []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `B`.*",
			wantJoin: "JOIN (`B`) ON (`A`.`B_id` = `B`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{Table: "B", Columns: selector{cols: []string{"c", "d"}}},
							Pairings:     []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.*, `B`.`c`, `B`.`d`",
			wantJoin: "JOIN (`B`) ON (`A`.`B_id` = `B`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{
								Table:   "B",
								Columns: selector{cols: []string{"c", "d"}},
							},
							Pairings: []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `B`.`c`, `B`.`d`",
			wantJoin: "JOIN (`B`) ON (`A`.`B_id` = `B`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{
								Table:   "B",
								Columns: selector{cols: []string{"c", "d"}},
							},
							Pairings: []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
						{
							SelectParams: common.SelectParams{
								Table:   "C",
								Columns: selector{cols: []string{"e", "f"}},
							},
							Pairings: []common.Pairing{{Column: "C_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `B`.`c`, `B`.`d`, `C`.`e`, `C`.`f`",
			wantJoin: "JOIN (`B`, `C`) ON (`A`.`B_id` = `B`.`id` AND `A`.`C_id` = `C`.`id`)",
		},
		{
			p: common.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []common.JoinParams{
						{
							SelectParams: common.SelectParams{
								Table: "B",
								Columns: selector{
									cols: []string{"c", "d"},
									joins: []common.JoinParams{
										{
											SelectParams: common.SelectParams{
												Table:   "D",
												Columns: selector{cols: []string{"g", "h"}},
											},
											Pairings: []common.Pairing{{Column: "D_id", JoinedColumn: "id"}},
										},
									},
								},
							},
							Pairings: []common.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
						{
							SelectParams: common.SelectParams{
								Table:   "C",
								Columns: selector{cols: []string{"e", "f"}},
							},
							Pairings: []common.Pairing{{Column: "C_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `B`.`c`, `B`.`d`, `D`.`g`, `D`.`h`, `C`.`e`, `C`.`f`",
			wantJoin: "JOIN (`B`, `C`) ON (`A`.`B_id` = `B`.`id` AND `A`.`C_id` = `C`.`id`) JOIN (`D`) ON (`B`.`D_id` = `D`.`id`)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantCols, func(t *testing.T) {
			assert.Equal(t, tt.wantCols, Columns(&tt.p))
			assert.Equal(t, tt.wantJoin, strings.Trim(Join(&tt.p), " "))
		})
	}
}

type selector struct {
	cols  []string
	joins []common.JoinParams
	count bool
}

func (s selector) Columns() []string {
	return s.cols
}

func (s selector) Joins() []common.JoinParams {
	return s.joins
}

func (s selector) Count() bool {
	return s.count
}
