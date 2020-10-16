package samples

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jf-tech/omniparser"
	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparserlegacy/omniv20"
	v20 "github.com/jf-tech/omniparserlegacy/omniv20/customfuncs"
)

// SampleTestCommon is a test helper for sample tests
func SampleTestCommon(t *testing.T, schemaFile, inputFile string) string {
	schemaFileBaseName := filepath.Base(schemaFile)
	schemaFileReader, err := os.Open(schemaFile)
	assert.NoError(t, err)
	defer schemaFileReader.Close()

	inputFileBaseName := filepath.Base(inputFile)
	inputFileReader, err := os.Open(inputFile)
	assert.NoError(t, err)
	defer inputFileReader.Close()

	schema, err := omniparser.NewSchema(
		schemaFileBaseName,
		schemaFileReader,
		omniparser.Extension{
			CreateSchemaHandler: omniv20.CreateSchemaHandler,
			CustomFuncs:         v20.OmniV20CustomFuncs,
		})
	assert.NoError(t, err)
	transform, err := schema.NewTransform(inputFileBaseName, inputFileReader, &transformctx.Ctx{})
	assert.NoError(t, err)

	var records []string
	for transform.Next() {
		recordBytes, err := transform.Read()
		assert.NoError(t, err)
		records = append(records, string(recordBytes))
	}

	return "[" + strings.Join(records, ",") + "]"
}
