package fixedlength

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"

	"github.com/jf-tech/omniparserlegacy/omniv20/samples"
)

func TestFixedLengthSample1_Simple(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./fixedlength_sample1_simple_schema.json", "./fixedlength_sample1_simple.txt", nil)))
}

func TestFixedLengthSample2_ByRows(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./fixedlength_sample2_by_rows_schema.json", "./fixedlength_sample2_by_rows.txt", nil)))
}

func TestFixedLengthSample3_ByHeaderFooter(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(
		samples.SampleTestCommon(
			t, "./fixedlength_sample3_by_header_footer_schema.json", "./fixedlength_sample3_by_header_footer.txt", nil)))
}

var bench = samples.NewBench(
	"./fixedlength_sample3_by_header_footer_schema.json", "./fixedlength_sample3_by_header_footer.txt", nil)

// BenchmarkFixedLengthSample3_ByHeaderFooter-8   	    3014	    366522 ns/op	   79245 B/op	    2248 allocs/op
func BenchmarkFixedLengthSample3_ByHeaderFooter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench.RunOneIteration(b)
	}
}
