package down

import (
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

//TODO this assumes people are using create or alter procedure (if the database supports it) 
func GetProcDiff(original map[string]common.Procedure, builder *strings.Builder){
	for _, proc := range original {
		builder.WriteString(proc.Definition+"\n")
	}
}
