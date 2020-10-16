package customfuncs

import (
	"sort"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"
)

func TestDumpOmniV20OnlyCustomFuncNames(t *testing.T) {
	var names []string
	for name := range OmniV20OnlyCustomFuncs {
		names = append(names, name)
	}
	sort.Strings(names)
	cupaloy.SnapshotT(t, jsons.BPM(names))
}
