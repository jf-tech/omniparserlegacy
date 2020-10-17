package omniv20

import (
	"fmt"
	"io"

	"github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/schemahandler"
	"github.com/jf-tech/omniparser/validation"

	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat"
	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat/delimited"
	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat/json"
	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat/xml"
	"github.com/jf-tech/omniparserlegacy/omniv20/transform"
	v20 "github.com/jf-tech/omniparserlegacy/omniv20/validation"
)

const (
	version = "omni.2.0"
)

// CreateSchemaHandler parses, validates and creates an omni-schema based handler.
func CreateSchemaHandler(ctx *schemahandler.CreateCtx) (schemahandler.SchemaHandler, error) {
	if ctx.Header.ParserSettings.Version != version {
		return nil, errs.ErrSchemaNotSupported
	}
	// First do a `transform_declarations` json schema validation
	err := validation.SchemaValidate(ctx.Name, ctx.Content, v20.JSONSchemaTransformDeclarations)
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

func fileFormats(ctx *schemahandler.CreateCtx) []fileformat.FileFormat {
	return []fileformat.FileFormat{
		delimited.NewDelimitedFileFormat(ctx.Name),
		json.NewJSONFileFormat(ctx.Name),
		xml.NewXMLFileFormat(ctx.Name),
		// TODO more built-in omni.2.0 file formats to come.
	}
}

type schemaHandler struct {
	ctx             *schemahandler.CreateCtx
	fileFormat      fileformat.FileFormat
	formatRuntime   interface{}
	finalOutputDecl *transform.Decl
}

func (h *schemaHandler) NewIngester(ctx *transformctx.Ctx, input io.Reader) (schemahandler.Ingester, error) {
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
