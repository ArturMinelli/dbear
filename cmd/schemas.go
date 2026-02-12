package cmd

import "strings"

// parseCommaSeparatedSchemas returns nil when flagValue is empty (meaning "all schemas"),
// otherwise returns trimmed non-empty schema names.
func parseCommaSeparatedSchemas(flagValue string) []string {
	trimmed := strings.TrimSpace(flagValue)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	schemas := make([]string, 0, len(parts))
	for _, part := range parts {
		schema := strings.TrimSpace(part)
		if schema != "" {
			schemas = append(schemas, schema)
		}
	}
	if len(schemas) == 0 {
		return nil
	}
	return schemas
}
