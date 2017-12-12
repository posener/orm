package dialect

import (
	"strings"
	"testing"

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
			wantCols: "`A`.*, `A_B_id`.*",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
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
			wantCols: "`A_B_id`.*, COUNT(*)",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
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
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.*",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
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
			wantCols: "`A`.*, `A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
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
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
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
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "JOIN (`B` AS `A_B_id`, `C` AS `A_C_id`) ON (`A`.`B_id` = `A_B_id`.`id` AND `A`.`C_id` = `A_C_id`.`id`)",
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
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_B_id_D_id`.`g`, `A_B_id_D_id`.`h`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "JOIN (`B` AS `A_B_id`, `C` AS `A_C_id`) ON (`A`.`B_id` = `A_B_id`.`id` AND `A`.`C_id` = `A_C_id`.`id`) JOIN (`D` AS `A_B_id_D_id`) ON (`A_B_id`.`D_id` = `A_B_id_D_id`.`id`)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantCols, func(t *testing.T) {
			d := &dialect{"mysql"}
			assert.Equal(t, tt.wantCols, columns(&tt.p))
			assert.Equal(t, tt.wantJoin, strings.Trim(d.join(&tt.p), " "))
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
