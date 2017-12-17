package dialect

import (
	"strings"
	"testing"

	"github.com/posener/orm/runtime"
	"github.com/stretchr/testify/assert"
)

func TestColumnsJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		p        runtime.SelectParams
		wantCols string
		wantJoin string
	}{
		{
			p: runtime.SelectParams{Table: "table", Columns: selector{}},
		},
		{
			p:        runtime.SelectParams{Table: "table", Columns: selector{count: true}},
			wantCols: "COUNT(*)",
		},
		{
			p:        runtime.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}}},
			wantCols: "`table`.`a`, `table`.`b`",
		},
		{
			p:        runtime.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}, count: true}},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p:        runtime.SelectParams{Table: "table", Columns: selector{cols: []string{"a", "b"}, count: true}},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{
								Table:   "B",
								Columns: selector{},
							},
							Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					count: true,
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{Table: "B", Columns: selector{}},
							Pairings:     []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "COUNT(*)",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{Table: "B", Columns: selector{}},
							Pairings:     []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{Table: "B", Columns: selector{cols: []string{"c", "d"}}},
							Pairings:     []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{
								Table:   "B",
								Columns: selector{cols: []string{"c", "d"}},
							},
							Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{
								Table:   "B",
								Columns: selector{cols: []string{"c", "d"}},
							},
							Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
						{
							SelectParams: runtime.SelectParams{
								Table:   "C",
								Columns: selector{cols: []string{"e", "f"}},
							},
							Pairings: []runtime.Pairing{{Column: "C_id", JoinedColumn: "id"}},
						},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "JOIN (`B` AS `A_B_id`, `C` AS `A_C_id`) ON (`A`.`B_id` = `A_B_id`.`id` AND `A`.`C_id` = `A_C_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Columns: selector{
					cols: []string{"a", "b"},
					joins: []runtime.JoinParams{
						{
							SelectParams: runtime.SelectParams{
								Table: "B",
								Columns: selector{
									cols: []string{"c", "d"},
									joins: []runtime.JoinParams{
										{
											SelectParams: runtime.SelectParams{
												Table:   "D",
												Columns: selector{cols: []string{"g", "h"}},
											},
											Pairings: []runtime.Pairing{{Column: "D_id", JoinedColumn: "id"}},
										},
									},
								},
							},
							Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
						},
						{
							SelectParams: runtime.SelectParams{
								Table:   "C",
								Columns: selector{cols: []string{"e", "f"}},
							},
							Pairings: []runtime.Pairing{{Column: "C_id", JoinedColumn: "id"}},
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
			d := Get("mysql").(*dialect)
			assert.Equal(t, tt.wantCols, d.selectColumns(&tt.p))
			assert.Equal(t, tt.wantJoin, strings.Trim(d.join(&tt.p), " "))
		})
	}
}

type selector struct {
	cols  []string
	joins []runtime.JoinParams
	count bool
}

func (s selector) Columns() []string {
	return s.cols
}

func (s selector) Joins() []runtime.JoinParams {
	return s.joins
}

func (s selector) Count() bool {
	return s.count
}
