package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const requestPropertyBecameEnumId = "request-property-became-enum"

func RequestPropertyBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {

						if enumDiff := propertyDiff.EnumDiff; enumDiff == nil || !enumDiff.EnumAdded {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        requestPropertyBecameEnumId,
							Level:     ERR,
							Text:      fmt.Sprintf(config.i18n(requestPropertyBecameEnumId), ColorizedValue(propertyFullName(propertyPath, propertyName))),
							Operation: operation,
							Path:      path,
							Source:    source,
						})
					})
			}
		}
	}
	return result
}
