package edi

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jf-tech/go-corelib/caches"
	"github.com/jf-tech/go-corelib/strs"
	errs "github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/validation"

	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat"
	"github.com/jf-tech/omniparserlegacy/omniv20/transform"
	v20 "github.com/jf-tech/omniparserlegacy/omniv20/validation"
)

const (
	fileFormatEDI = "edi"
)

type ediFileFormat struct {
	schemaName string
}

// NewEDIFileFormat creates a FileFormat for EDI.
func NewEDIFileFormat(schemaName string) fileformat.FileFormat {
	return &ediFileFormat{schemaName: schemaName}
}

type ediFormatRuntime struct {
	Decl  *fileDecl `json:"file_declaration"`
	XPath string
}

func (f *ediFileFormat) ValidateSchema(
	format string, schemaContent []byte, finalOutputDecl *transform.Decl) (interface{}, error) {
	if format != fileFormatEDI {
		return nil, errs.ErrSchemaNotSupported
	}
	err := validation.SchemaValidate(f.schemaName, schemaContent, v20.JSONSchemaEDIFileDeclaration)
	if err != nil {
		// err is already context formatted.
		return nil, err
	}
	var runtime ediFormatRuntime
	_ = json.Unmarshal(schemaContent, &runtime) // JSON schema validation earlier guarantees Unmarshal success.
	err = f.validateFileDecl(runtime.Decl)
	if err != nil {
		// err is already context formatted.
		return nil, err
	}
	if finalOutputDecl == nil {
		return nil, f.FmtErr("'FINAL_OUTPUT' is missing")
	}
	runtime.XPath = strings.TrimSpace(strs.StrPtrOrElse(finalOutputDecl.XPath, ""))
	if runtime.XPath != "" {
		_, err := caches.GetXPathExpr(runtime.XPath)
		if err != nil {
			return nil, f.FmtErr("'FINAL_OUTPUT.xpath' (value: '%s') is invalid, err: %s",
				runtime.XPath, err.Error())
		}
	}
	return &runtime, nil
}

func (f *ediFileFormat) validateFileDecl(decl *fileDecl) error {
	err := (&ediValidateCtx{}).validateFileDecl(decl)
	if err != nil {
		return f.FmtErr(err.Error())
	}
	return err
}

func (f *ediFileFormat) CreateFormatReader(
	name string, r io.Reader, runtime interface{}) (fileformat.FormatReader, error) {
	edi := runtime.(*ediFormatRuntime)
	return NewReader(name, r, edi.Decl, edi.XPath)
}

func (f *ediFileFormat) FmtErr(format string, args ...interface{}) error {
	return fmt.Errorf("schema '%s': %s", f.schemaName, fmt.Sprintf(format, args...))
}
