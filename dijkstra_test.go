package dijkstra

import (
	"reflect"
	"testing"
)

func TestNoPath(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, true)
}

func TestLoop(t *testing.T) {
	testSolution(t, BestPath{}, newErrLoop(1, 2), "testdata/J.txt", 0, 4, true)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true)
	testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false)
	testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true)
}
func BenchmarkDmitrySigaevNodes4(b *testing.B)    { benchmarkAlt(b, "testdata/4.txt", 0) }
func BenchmarkDmitrySigaevNodes10(b *testing.B)   { benchmarkAlt(b, "testdata/10.txt", 0) }
func BenchmarkDmitrySigaevNodes100(b *testing.B)  { benchmarkAlt(b, "testdata/100.txt", 0) }
func BenchmarkDmitrySigaevNodes1000(b *testing.B) { benchmarkAlt(b, "testdata/1000.txt", 0) }

func benchmarkAlt(b *testing.B, filename string, i int) {
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

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int, shortest bool) {
	var err error
	var graph Graph
	graph, err = Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	if shortest {
		got, err = graph.Shortest(from, to)
	} else {
		got, err = graph.Longest(from, to)
	}
	testErrors(t, wanterr, err, filename)
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
