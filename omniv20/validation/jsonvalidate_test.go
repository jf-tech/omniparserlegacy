package validation

import (
	"testing"

	"github.com/jf-tech/omniparser/validation"
	"github.com/stretchr/testify/assert"
)

func TestSchemaValidate(t *testing.T) {
	for _, test := range []struct {
		name          string
		jsonSchema    string
		schemaContent string
		expectedErr   string
	}{
		{
			name:       "success",
			jsonSchema: JSONSchemaCSVFileDeclaration,
			schemaContent: `{
					"file_declaration": {
						"delimiter": "|",
						"replace_double_quotes": true,
						"header_row_index": 2,
						"data_row_index": 4,
						"columns": [
							{ "name": "COL1" },
							{ "name": "COL 2", "alias": "COL2" }
						]
					}
				}`,
			expectedErr: "",
		},
		{
			name:          "invalid json schema",
			jsonSchema:    ">>",
			schemaContent: `{}`,
			expectedErr:   `unable to perform schema validation: invalid character '>' looking for beginning of value`,
		},
		{
			name:       "invalid delimiter",
			jsonSchema: JSONSchemaCSVFileDeclaration,
			schemaContent: `{
					"file_declaration": {
						"delimiter": "length longer than 1",
						"data_row_index": 4,
						"columns": [ { "name": "COL1" } ]
					}
				}`,
			expectedErr: `schema 'test-schema' validation failed: file_declaration.delimiter: String length must be less than or equal to 1`,
		},
		{
			name:       "multiple errors",
			jsonSchema: JSONSchemaCSVFileDeclaration,
			schemaContent: `{
					"file_declaration": {
						"delimiter": "length longer than 1",
						"data_row_index": -1,
						"columns": [ { "name": "COL1" } ]
					}
				}`,
			expectedErr: "schema 'test-schema' validation failed:\nfile_declaration.data_row_index: Must be greater than or equal to 1\nfile_declaration.delimiter: String length must be less than or equal to 1",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validation.SchemaValidate("test-schema", []byte(test.schemaContent), test.jsonSchema)
			if test.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
