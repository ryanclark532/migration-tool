package down

import (
	"fmt"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

// TODO this assumes people are using create or alter procedure (if the database supports it)
func GetProcDiff(original map[string]common.Procedure, post map[string]common.Procedure, builder *strings.Builder) {
	//If exists in post but not pre drop procedure
	//If exists in pre but not post create procedure
	//If exists in both create "ALTER" procedure

	processedProcs := make(map[string]bool)

	for key, _ := range original {
		if prev, ex := post[key]; ex {
			builder.WriteString(prev.Definition)
		}
		processedProcs[key] = true
	}

	for key, _ := range post {
		if _, exists := processedProcs[key]; !exists {
			builder.WriteString(fmt.Sprintf("DROP PROCEDURE %s", key))
		}
	}
}
