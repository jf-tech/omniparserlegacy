package validation

//go:generate sh -c "cd jsonschemas && go run gen.go -json transform_declarations.json -varname JSONSchemaTransformDeclarations > ../transformDeclarations.go"
//go:generate sh -c "cd jsonschemas && go run gen.go -json csv_file_declaration.json -varname JSONSchemaCSVFileDeclaration > ../csvFileDeclaration.go"
