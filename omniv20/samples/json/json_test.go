package json

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"

	"github.com/jf-tech/omniparserlegacy/omniv20/samples"
)

func TestJSONSample1(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./json_sample1_schema.json", "./json_sample1.json", nil)))
}

func TestJSONSample2XPathDynamic(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./json_sample2_xpathdynamic_schema.json", "./json_sample2_xpathdynamic.json", nil)))
}

var bench = samples.NewBench("./json_sample1_schema.json", "./json_sample1.json", nil)

// go test -bench=. -benchmem
// BenchmarkJSONSample1-8   	    4821	    230071 ns/op	   72929 B/op	    1927 allocs/op

func BenchmarkJSONSample1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench.RunOneIteration(b)
	}
}
