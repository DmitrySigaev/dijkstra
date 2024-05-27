package dijkstra

import (
	"reflect"
	"testing"
)

func TestFailure(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5)
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
	graph, _, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Shortest(src, dest)
	}
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int) {
	graph, _, err := Import(filename)
	if err != nil {
		t.Error(err)
	}
	got, err := graph.Shortest(from, to)
	testErrors(t, wanterr, err)
	if got.Distance != best.Distance {
		t.Error("Distance incorrect\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error("Path incorrect\ngot: ", got.Path, "\nwant: ", best.Path)
	}
}
