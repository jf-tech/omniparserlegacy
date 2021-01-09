package omniv20

import (
	"encoding/json"
	"errors"

	"github.com/jf-tech/omniparser/customfuncs"
	"github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat"
	"github.com/jf-tech/omniparserlegacy/omniv20/transform"
)

type ingester struct {
	finalOutputDecl *transform.Decl
	customFuncs     customfuncs.CustomFuncs
	ctx             *transformctx.Ctx
	reader          fileformat.FormatReader
}

func (g *ingester) Read() (interface{}, []byte, error) {
	n, err := g.reader.Read()
	if err != nil {
		// Read() supposed to have already done CtxAwareErr error wrapping. So directly return.
		return nil, nil, err
	}
	defer g.reader.Release(n)
	result, err := transform.NewParseCtx(g.ctx, g.customFuncs).ParseNode(n, g.finalOutputDecl)
	if err != nil {
		// ParseNode() error not CtxAwareErr wrapped, so wrap it.
		// Note errs.ErrorTransformFailed is a continuable error.
		return nil, nil, errs.ErrTransformFailed(g.fmtErrStr("fail to transform. err: %s", err.Error()))
	}
	transformed, err := json.Marshal(result)
	return nil, transformed, err
}

func (g *ingester) IsContinuableError(err error) bool {
	return errs.IsErrTransformFailed(err) || g.reader.IsContinuableError(err)
}

func (g *ingester) FmtErr(format string, args ...interface{}) error {
	return errors.New(g.fmtErrStr(format, args...))
}

func (g *ingester) fmtErrStr(format string, args ...interface{}) string {
	return g.reader.FmtErr(format, args...).Error()
}
