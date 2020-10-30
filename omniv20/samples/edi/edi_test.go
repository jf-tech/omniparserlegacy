package edi

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"

	"github.com/jf-tech/omniparserlegacy/omniv20/samples"
)

func TestEDISample1(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./edi_sample1_schema.json", "./edi_sample1.txt", nil)))
}

var bench = samples.NewBench("./edi_sample1_schema.json", "./edi_sample1.txt", nil)

// BenchmarkEDISample1-8   	  148292	    240290 ns/op	   78183 B/op	    1622 allocs/op
func BenchmarkEDISample1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench.RunOneIteration(b)
	}
}
