package common

var MySQLTypeMap = map[string]string{
	"tinyint":            "int8",
	"smallint":           "int16",
	"mediumint":          "int32",
	"int":                "int32",
	"bigint":             "int64",
	"tinyint unsigned":   "uint8",
	"smallint unsigned":  "uint16",
	"mediumint unsigned": "uint",
	"int unsigned":       "uint32",
	"bigint unsigned":    "uint64",

	"float":              "float32",
	"double":             "float64",
	"decimal":            "float64",

	"char":       "string",
	"varchar":    "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"json":       "string",

	"date":      "[]uint8",
	"datetime":  "[]uint8",
	"time":      "[]uint8",
	"timestamp": "[]uint8",

	"enum": "string",

	"binary":     "[]byte",
	"varbinary":  "[]byte",
	"tinyblob":   "[]byte",
	"blob":       "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",
}
