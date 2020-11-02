package samples

import (
	"bytes"
	"io"
	"io/ioutil"
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
func SampleTestCommon(t *testing.T, schemaFile, inputFile string, externals map[string]string) string {
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
	transform, err := schema.NewTransform(
		inputFileBaseName,
		inputFileReader,
		&transformctx.Ctx{ExternalProperties: externals})
	assert.NoError(t, err)

	var records []string
	for {
		recordBytes, err := transform.Read()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		records = append(records, string(recordBytes))
	}

	return "[" + strings.Join(records, ",") + "]"
}

type Bench struct {
	schema    omniparser.Schema
	input     []byte
	externals map[string]string
}

func (bch *Bench) RunOneIteration(b *testing.B) {
	transform, err := bch.schema.NewTransform(
		"bench",
		bytes.NewReader(bch.input),
		&transformctx.Ctx{ExternalProperties: bch.externals})
	if err != nil {
		b.FailNow()
	}
	for {
		_, err := transform.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			b.FailNow()
		}
	}
}

func NewBench(schemaFile, inputFile string, externals map[string]string) *Bench {
	schemaContent, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		panic(err)
	}
	schema, err := omniparser.NewSchema("bench", bytes.NewReader(schemaContent),
		omniparser.Extension{
			CreateSchemaHandler: omniv20.CreateSchemaHandler,
			CustomFuncs:         v20.OmniV20CustomFuncs,
		})
	if err != nil {
		panic(err)
	}
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	return &Bench{
		schema:    schema,
		input:     input,
		externals: externals,
	}
}
