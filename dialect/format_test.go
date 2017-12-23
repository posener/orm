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
			p: runtime.SelectParams{Table: "table"},
		},
		{
			p:        runtime.SelectParams{Table: "table", Count: true},
			wantCols: "COUNT(*)",
		},
		{
			p: runtime.SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
			},
			wantCols: "`table`.`a`, `table`.`b`",
		},
		{
			p: runtime.SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Count:          true,
			},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: runtime.SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"b", "a"},
				Count:          true,
			},
			wantCols: "`table`.`b`, `table`.`a`, COUNT(*)",
		},
		{
			p: runtime.SelectParams{
				Table:          "table",
				Columns:        map[string]bool{"a": true, "b": true},
				Count:          true,
				OrderedColumns: []string{"a", "b"},
			},
			wantCols: "`table`.`a`, `table`.`b`, COUNT(*)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{
							Table: "B",
						},
						Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Count: true,
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{Table: "B"},
						Pairings:     []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "COUNT(*)",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{Table: "B"},
						Pairings:     []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table: "A",
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`) ON (`A`.`B_id` = `A_B_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
						},
						Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
					{
						SelectParams: runtime.SelectParams{
							Table:          "C",
							Columns:        map[string]bool{"e": true, "f": true},
							OrderedColumns: []string{"e", "f"},
						},
						Pairings: []runtime.Pairing{{Column: "C_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`, `C` AS `A_C_id`) ON (`A`.`B_id` = `A_B_id`.`id` AND `A`.`C_id` = `A_C_id`.`id`)",
		},
		{
			p: runtime.SelectParams{
				Table:          "A",
				Columns:        map[string]bool{"a": true, "b": true},
				OrderedColumns: []string{"a", "b"},
				Joins: []runtime.JoinParams{
					{
						SelectParams: runtime.SelectParams{
							Table:          "B",
							Columns:        map[string]bool{"c": true, "d": true},
							OrderedColumns: []string{"c", "d"},
							Joins: []runtime.JoinParams{
								{
									SelectParams: runtime.SelectParams{
										Table:          "D",
										Columns:        map[string]bool{"g": true, "h": true},
										OrderedColumns: []string{"g", "h"},
									},
									Pairings: []runtime.Pairing{{Column: "D_id", JoinedColumn: "id"}},
								},
							},
						},
						Pairings: []runtime.Pairing{{Column: "B_id", JoinedColumn: "id"}},
					},
					{
						SelectParams: runtime.SelectParams{
							Table:          "C",
							Columns:        map[string]bool{"e": true, "f": true},
							OrderedColumns: []string{"e", "f"},
						},
						Pairings: []runtime.Pairing{{Column: "C_id", JoinedColumn: "id"}},
					},
				},
			},
			wantCols: "`A`.`a`, `A`.`b`, `A_B_id`.`c`, `A_B_id`.`d`, `A_B_id_D_id`.`g`, `A_B_id_D_id`.`h`, `A_C_id`.`e`, `A_C_id`.`f`",
			wantJoin: "LEFT OUTER JOIN (`B` AS `A_B_id`, `C` AS `A_C_id`) ON (`A`.`B_id` = `A_B_id`.`id` AND `A`.`C_id` = `A_C_id`.`id`) LEFT OUTER JOIN (`D` AS `A_B_id_D_id`) ON (`A_B_id`.`D_id` = `A_B_id_D_id`.`id`)",
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
