package codegenerator

import (
	"testing"
)

func TestCodeGenerator(t *testing.T) {

	code("t_user")

}
func TestCodeGeneratorALL(t *testing.T) {
	tableNames := selectAllTable()
	for _, tableName := range tableNames {
		code(tableName)
	}
}
