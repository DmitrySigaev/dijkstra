package dijkstra

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestNoPath(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, true, -1)
}

func TestLoop(t *testing.T) {
	testSolution(t, BestPath{}, newErrLoop(2, 1), "testdata/J.txt", 0, 4, true, -1)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true, -1)
	testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false, -1)
	testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true, -1)
}

func TestCorrectAllLists(t *testing.T) {
	for i := 0; i <= 3; i++ {
		testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true, i)
		testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false, i)
		testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true, i)
	}
}

func TestCorrectAutoLargeList(t *testing.T) {
	g := NewGraph()
	for i := 0; i < 2000; i++ {
		v := g.AddNewVertex()
		v.AddArc(i+1, 1)
	}
	g.AddNewVertex()
	_, err := g.Shortest(0, 2000)
	testErrors(t, nil, err, "manual test")
	_, err = g.Longest(0, 2000)
	testErrors(t, nil, err, "manual test")
}

var benchNames = []string{"github.com/DmitrySigaev"}
var listNames = []string{"PQShort", "PQLong", "LLShort", "LLLong"}

func BenchmarkSetup(b *testing.B) {
	nodeIterations := 6
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		b.Run("setup/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
			filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
			if _, err := os.Stat(filename); err != nil {
				g := Generate(nodes)
				err := g.ExportToFile(filename)
				if err != nil {
					log.Fatal(err)
				}
			}
			g, _ := Import(filename)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g.setup(true, 0, -1)
			}
		})
	}
}

func BenchmarkLists(b *testing.B) {
	nodeIterations := 6
	shortest := false
	shortText := []string{"Short", "Long"}
	for z := 0; z < 2; z++ {
		shortest = !shortest
		for i, n := range listNames {
			nodes := 1
			for j := 0; j < nodeIterations; j++ {
				nodes *= 4
				b.Run(shortText[z]+"/"+n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
					benchmarkList(b, nodes, i, shortest)
				})
			}
		}
	}
}

func benchmarkList(b *testing.B, nodes, list int, shortest bool) {

	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.setup(shortest, src, list)
		graph.postSetupEvaluate(src, dest, shortest)
	}
}

func BenchmarkAll(b *testing.B) {
	nodeIterations := 6
	for i, n := range benchNames {
		nodes := 1
		for j := 0; j < nodeIterations; j++ {
			nodes *= 4
			b.Run(n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
				benchmarkAlt(b, nodes, i)
			})

		}
	}
	//Cleanup
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		os.Remove("testdata/bench/" + strconv.Itoa(nodes) + ".txt")
	}
}

/*
//Mattomatics does not work.
func BenchmarkMattomaticNodes4(b *testing.B)    { benchmarkAlt(b, 4, 3) }
*/
func benchmarkAlt(b *testing.B, nodes, i int) {
	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch i {
	case 0:
		benchmarkRC(b, filename)
	default:
		b.Error("You're retarded")
	}
}

func benchmarkRC(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Shortest(src, dest)
	}
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int, shortest bool, list int) {
	var err error
	var graph Graph
	graph, err = Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	if list >= 0 {
		graph.setup(shortest, from, list)
		got, err = graph.postSetupEvaluate(from, to, shortest)
	} else if shortest {
		got, err = graph.Shortest(from, to)
	} else {
		got, err = graph.Longest(from, to)
	}
	testErrors(t, wanterr, err, filename)
	testResults(t, got, best, shortest, filename)
}

func testResults(t *testing.T, got, best BestPath, shortest bool, filename string) {
	distmethod := "Shortest"
	if !shortest {
		distmethod = "Longest"
	}
	if got.Distance != best.Distance {
		t.Error(distmethod, " distance incorrect\n", filename, "\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error(distmethod, " path incorrect\n\n", filename, "got: ", got.Path, "\nwant: ", best.Path)
	}
}

func getKSolLong() BestPath {
	return BestPath{
		31,
		[]int{
			0, 1, 3, 2, 4,
		},
	}
}
func getKSolShort() BestPath {
	return BestPath{
		2,
		[]int{
			0, 3, 4,
		},
	}
}
