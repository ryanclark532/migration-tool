package down

import (
	"fmt"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

func GetProcDiff(original map[string]common.Procedure, post map[string]common.Procedure, builder *strings.Builder, processedProcs map[string]bool) {
	//If exists in post but not pre drop procedure
	//If exists in pre but not post create procedure
	//If exists in both create "ALTER" procedure

	for key, ori := range original {
		if _, exists := processedProcs[key]; exists {
			continue
		}
		if prev, ex := post[key]; ex {
			if prev.Definition == ori.Definition {
				builder.WriteString(prev.Definition)
			}
		}
		processedProcs[key] = true
	}

	for key, _ := range post {
		if _, exists := processedProcs[key]; !exists {
			builder.WriteString(fmt.Sprintf("DROP PROCEDURE %s", key))
		}
	}
}
