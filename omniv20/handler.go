package omniv20

import (
	"fmt"
	"io"

	"github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/handlers"

	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat"
	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat/delimited"
	"github.com/jf-tech/omniparserlegacy/omniv20/transform"
	"github.com/jf-tech/omniparserlegacy/omniv20/validation"
)

const (
	version = "omni.2.0"
)

// CreateHandler parses, validates and creates an omni-schema based handler.
func CreateHandler(ctx *handlers.HandlerCtx) (handlers.SchemaHandler, error) {
	if ctx.Header.ParserSettings.Version != version {
		return nil, errs.ErrSchemaNotSupported
	}
	// First do a `transform_declarations` json schema validation
	err := validation.SchemaValidate(ctx.Name, ctx.Content, validation.JSONSchemaTransformDeclarations)
	if err != nil {
		// err is already context formatted.
		return nil, err
	}
	finalOutputDecl, err := transform.ValidateTransformDeclarations(ctx.Content, ctx.CustomFuncs)
	if err != nil {
		return nil, fmt.Errorf(
			"schema '%s' 'transform_declarations' validation failed: %s",
			ctx.Name, err.Error())
	}
	for _, format := range fileFormats(ctx) {
		formatRuntime, err := format.ValidateSchema(
			ctx.Header.ParserSettings.FileFormatType,
			ctx.Content,
			finalOutputDecl)
		if err == errs.ErrSchemaNotSupported {
			continue
		}
		if err != nil {
			// error from FileFormat is already context formatted.
			return nil, err
		}
		return &schemaHandler{
			ctx:             ctx,
			fileFormat:      format,
			formatRuntime:   formatRuntime,
			finalOutputDecl: finalOutputDecl,
		}, nil
	}
	return nil, errs.ErrSchemaNotSupported
}

func fileFormats(ctx *handlers.HandlerCtx) []fileformat.FileFormat {
	return []fileformat.FileFormat{
		delimited.NewDelimitedFileFormat(ctx.Name),
		// TODO more built-in omni.2.0 file formats to come.
	}
}

type schemaHandler struct {
	ctx             *handlers.HandlerCtx
	fileFormat      fileformat.FileFormat
	formatRuntime   interface{}
	finalOutputDecl *transform.Decl
}

func (h *schemaHandler) NewIngester(ctx *transformctx.Ctx, input io.Reader) (handlers.Ingester, error) {
	reader, err := h.fileFormat.CreateFormatReader(ctx.InputName, input, h.formatRuntime)
	if err != nil {
		return nil, err
	}
	return &ingester{
		finalOutputDecl: h.finalOutputDecl,
		customFuncs:     h.ctx.CustomFuncs,
		ctx:             ctx,
		reader:          reader,
	}, nil
}
