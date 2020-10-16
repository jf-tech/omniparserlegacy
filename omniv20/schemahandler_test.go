package omniv20

import (
	"errors"
	"io"
	"testing"

	"github.com/jf-tech/omniparser/idr"
	"github.com/jf-tech/omniparser/schemahandler"
	"github.com/jf-tech/omniparser/transformctx"
	"github.com/stretchr/testify/assert"

	"github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/header"

	"github.com/jf-tech/omniparserlegacy/omniv20/fileformat"
	"github.com/jf-tech/omniparserlegacy/omniv20/transform"
)

func TestCreateHandler_VersionNotSupported(t *testing.T) {
	p, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version: "12345",
				},
			},
		})
	assert.Error(t, err)
	assert.Equal(t, errs.ErrSchemaNotSupported, err)
	assert.Nil(t, p)
}

func TestCreateHandler_FormatNotSupported(t *testing.T) {
	p, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version:        version,
					FileFormatType: "unknown",
				},
			},
			Content: []byte(`{"transform_declarations": { "FINAL_OUTPUT": {} }}`),
		})
	assert.Error(t, err)
	assert.Equal(t, errs.ErrSchemaNotSupported, err)
	assert.Nil(t, p)
}

func TestCreateHandler_TransformDeclarationsJSONValidationFailed(t *testing.T) {
	p, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Name: "test-schema",
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version:        version,
					FileFormatType: "xml",
				},
			},
			Content: []byte(`{"transform_declarations": {}}`),
		})
	assert.Error(t, err)
	assert.Equal(t,
		`schema 'test-schema' validation failed: transform_declarations: FINAL_OUTPUT is required`,
		err.Error())
	assert.Nil(t, p)
}

func TestCreateHandler_TransformDeclarationsInCodeValidationFailed(t *testing.T) {
	p, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Name: "test-schema",
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version:        version,
					FileFormatType: "xml",
				},
			},
			Content: []byte(
				`{
					"transform_declarations": {
						"FINAL_OUTPUT": { "template": "non-existing" }
					}
				}`),
		})
	assert.Error(t, err)
	assert.Equal(t,
		`schema 'test-schema' 'transform_declarations' validation failed: 'FINAL_OUTPUT' contains non-existing template reference 'non-existing'`,
		err.Error())
	assert.Nil(t, p)
}

func TestCreateHandler_FileFormatValidationFailed(t *testing.T) {
	p, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Name: "test-schema",
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version:        version,
					FileFormatType: "delimited",
				},
			},
			Content: []byte(
				`{
					"file_declaration": {
						"delimiter": ",",
						"data_row_index": -1,
						"columns": [ { "name": "col1" } ]
					},	
					"transform_declarations": {
						"FINAL_OUTPUT": { "xpath": "." }
					}
				}`),
		})
	assert.Error(t, err)
	assert.Equal(t,
		`schema 'test-schema' validation failed: file_declaration.data_row_index: Must be greater than or equal to 1`,
		err.Error())
	assert.Nil(t, p)
}

func TestCreateHandler_Success(t *testing.T) {
	h, err := CreateSchemaHandler(
		&schemahandler.CreateCtx{
			Name: "test-schema",
			Header: header.Header{
				ParserSettings: header.ParserSettings{
					Version:        version,
					FileFormatType: "delimited",
				},
			},
			Content: []byte(
				`{
					"file_declaration": {
						"delimiter": ",",
						"data_row_index": 1,
						"columns": [ { "name": "col1" } ]
					},	
					"transform_declarations": {
						"FINAL_OUTPUT": { "xpath": "." }
					}
				}`),
		})
	assert.NoError(t, err)
	ingester, err := h.NewIngester(&transformctx.Ctx{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ingester)
}

type testFileFormat struct {
	createFormatReaderErr error
}

func (f testFileFormat) ValidateSchema(_ string, _ []byte, _ *transform.Decl) (interface{}, error) {
	return nil, nil
}

func (f testFileFormat) CreateFormatReader(
	inputName string, input io.Reader, runtime interface{}) (fileformat.FormatReader, error) {
	if f.createFormatReaderErr != nil {
		return nil, f.createFormatReaderErr
	}
	return testFormatReader{
		inputName: inputName,
		input:     input,
		runtime:   runtime,
	}, nil
}

type testFormatReader struct {
	inputName string
	input     io.Reader
	runtime   interface{}
}

func (r testFormatReader) Read() (*idr.Node, error)            { panic("implement me") }
func (r testFormatReader) Release(*idr.Node)                   { panic("implement me") }
func (r testFormatReader) IsContinuableError(error) bool       { panic("implement me") }
func (r testFormatReader) FmtErr(string, ...interface{}) error { panic("implement me") }

func TestNewIngester_Failure(t *testing.T) {
	h := &schemaHandler{
		ctx:        &schemahandler.CreateCtx{},
		fileFormat: &testFileFormat{createFormatReaderErr: errors.New("create reader failure")},
	}
	ingester, err := h.NewIngester(&transformctx.Ctx{}, nil)
	assert.Error(t, err)
	assert.Equal(t, "create reader failure", err.Error())
	assert.Nil(t, ingester)
}
