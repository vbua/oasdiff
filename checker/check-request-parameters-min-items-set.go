package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinItemsSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}
					minItemsDiff := paramDiff.SchemaDiff.MinItemsDiff
					if minItemsDiff == nil {
						continue
					}
					if minItemsDiff.From != nil ||
						minItemsDiff.To == nil {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]
					result = append(result, BackwardCompatibilityError{
						Id:        "request-parameter-min-items-set",
						Level:     WARN,
						Text:      fmt.Sprintf("for the %s request parameter %s, the minItems was set to '%s'", ColorizedValue(paramLocation), ColorizedValue(paramName), minItemsDiff.To),
						Comment:   "It is warn because sometimes it is required to be set because of security reasons or current error in specification. But good clients should be checked to support this restriction before such change in specification.",
						Operation: operation,
						Path:      path,
						Source:    source,
						ToDo:      "Add to exceptions-list.md",
					})
				}
			}
		}
	}
	return result
}
