package xml

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"

	"github.com/jf-tech/omniparserlegacy/omniv20/samples"
)

// TODO: this sample is slightly modified from the original to remove the not-yet-supported xpath func 'sort'.
func TestXMLSample1(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./xml_sample1_schema.json", "./xml_sample1.xml", nil)))
}

var bench = samples.NewBench("./xml_sample1_schema.json", "./xml_sample1.xml", nil)

// go test -bench=. -benchmem
// BenchmarkXMLSample1-8   	    4696	    233774 ns/op	   70624 B/op	    1785 allocs/op

func BenchmarkXMLSample1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench.RunOneIteration(b)
	}
}
