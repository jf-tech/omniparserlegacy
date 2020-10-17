package csv

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"

	"github.com/jf-tech/omniparserlegacy/omniv20/samples"
)

var csvSample1Externals = map[string]string{
	"source_info": "wiki",
	"location":    "india",
}

func TestCSVSample1(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./csv_sample1_schema.json", "./csv_sample1.csv", csvSample1Externals)))
}

var bench = samples.NewBench("./csv_sample1_schema.json", "./csv_sample1.csv", csvSample1Externals)

// go test -bench=. -benchmem -benchtime=30s
// BenchmarkCSVSample1-8   	  148292	    240290 ns/op	   78183 B/op	    1622 allocs/op

func BenchmarkCSVSample1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench.RunOneIteration(b)
	}
}
