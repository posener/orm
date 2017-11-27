package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	benchLine = regexp.MustCompile(`^Benchmark(ORM|GORM|Raw)([^-]+)-(\d+)\W+(\d+)\W+(\d+)\W+ns\/op$`)
	pkgs      = []string{"ORM", "Raw", "GORM"}
)

func main() {
	f, err := os.Open("benchmark.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tests, results := parse(f)

	groups := map[string]plotter.Values{}
	for pkg, values := range results {
		for i := range values {
			groups[pkg] = append(groups[pkg], float64(values[i]))
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Benchmark"
	p.Y.Label.Text = "ns/op"

	w := float64(20)
	i := 0
	for _, name := range pkgs {
		values := groups[name]
		bars, err := plotter.NewBarChart(values, vg.Points(w))
		if err != nil {
			panic(err)
		}
		bars.LineStyle.Width = vg.Length(1)
		bars.Color = plotutil.Color(i)
		bars.Offset = vg.Points(-w + float64(i)*w)
		p.Add(bars)
		p.Legend.Add(name, bars)
		i++
	}
	p.Legend.Top = true
	p.NominalX(tests...)

	if err := p.Save(5*vg.Inch, 3*vg.Inch, "benchmark.png"); err != nil {
		panic(err)
	}

	markdown(tests, results)
}

func parse(r io.Reader) (tests []string, results map[string][]int) {
	results = make(map[string][]int)
	b := bufio.NewReader(r)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		matches := benchLine.FindStringSubmatch(string(line))
		if len(matches) < 6 {
			continue
		}
		pkg := matches[1]
		test := matches[2]
		result := matches[5]

		testIndex := find(tests, test)
		if testIndex == -1 {
			tests = append(tests, test)
			testIndex = len(tests) - 1
		}

		i, err := strconv.Atoi(result)
		if err != nil {
			panic(err)
		}
		if len(results[pkg]) == 0 {
			results[pkg] = make([]int, 100)
		}
		results[pkg][testIndex] = i
	}

	for pkg := range results {
		results[pkg] = results[pkg][:len(tests)]
	}

	return tests, results
}

func markdown(tests []string, results map[string][]int) {
	f, err := os.Create("./README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("# Latest benchmark results\n\n")
	f.WriteString("Units: [ns/op (ratio from ORM)]\n\n")
	f.WriteString("| Test | " + strings.Join(pkgs, " | ") + " |\n")
	f.WriteString("| --- | --- | --- | --- |\n")
	for i := range tests {
		f.WriteString("| " + tests[i] + " | " + strings.Join(values(results, i), " | ") + " |\n")
	}
	f.WriteString("\n\n![graph](./benchmark.png)\n\n")
	f.WriteString(fmt.Sprintf("\n\nBenchmark time: %s\n", time.Now().UTC().Format("2006-01-02")))
	f.WriteString(`

#### Compared packages:

- [x] ORM: posener/orm (this package)
- [x] RAW: Direct SQL commands
- [x] GORM: jinzhu/gorm

#### Operations:

- [x] Insert: INSERT operations
- [X] Query: SELECT operations on an object with 2 fields
- [X] QueryLargeStruct: SELECT on an object of ~35 different fields

`)
}

func values(results map[string][]int, i int) []string {
	var ret []string
	for j := 0; j < len(pkgs); j++ {
		cur := results[pkgs[j]][i]
		percent := int(float64(cur) / float64(results[pkgs[0]][i]) * 100.0)
		ret = append(ret, fmt.Sprintf("%d (%d%%)", cur, percent))
	}
	return ret
}

func find(arr []string, s string) int {
	for i := range arr {
		if arr[i] == s {
			return i
		}
	}
	return -1
}
