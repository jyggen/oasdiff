package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const requestBodyBecameEnumId = "request-body-became-enum"

func RequestBodyBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
				if schemaDiff := mediaTypeDiff.SchemaDiff; schemaDiff.EnumDiff == nil || !schemaDiff.EnumDiff.EnumAdded {
					continue
				}
				result = append(result, BackwardCompatibilityError{
					Id:        requestBodyBecameEnumId,
					Level:     ERR,
					Text:      config.i18n(requestBodyBecameEnumId),
					Operation: operation,
					Path:      path,
					Source:    source,
				})
			}
		}
	}
	return result
}
