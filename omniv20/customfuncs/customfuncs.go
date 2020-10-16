package customfuncs

import (
	"github.com/jf-tech/omniparser/customfuncs"

	v10 "github.com/jf-tech/omniparserlegacy/omniv10/customfuncs"
)

// OmniV20OnlyCustomFuncs contains 'omni.2.0' specific custom funcs.
var OmniV20OnlyCustomFuncs = map[string]customfuncs.CustomFuncType{
	// keep these custom funcs lexically sorted
	"avg": Avg,
	"sum": Sum,
}

// OmniV20CustomFuncs all custom funcs supported by 'omni.2.0'.
var OmniV20CustomFuncs = customfuncs.Merge(
	customfuncs.CommonCustomFuncs,
	v10.OmniV10CustomFuncs,
	OmniV20OnlyCustomFuncs)
