package rules

func init_database() {
	Register("Exists", validateExists)
	Register("Unique", validateUnique)
}

// validateExists: the value must exist as a column value in the given table.
// Params: [table, column, extraKey, extraVal, ...]
// Requires a PresenceVerifier to be set on the validator; passes trivially if none.
func validateExists(_ string, value any, params []string, ctx RuleContext) bool {
	pv := ctx.GetPresenceVerifier()

	if pv == nil {
		return true // no verifier configured — skip
	}

	if len(params) < 2 {
		return true
	}

	table := params[0]
	column := params[1]

	extras := parseExtras(params[2:])

	return pv.GetCount(table, column, stringify(value), nil, nil, extras) >= 1
}

// validateUnique: the value must NOT already exist in the given table/column.
// Params: [table, column, ignoreValue, ignoreColumn, extraKey, extraVal, ...]
// Requires a PresenceVerifier; passes trivially if none.
func validateUnique(_ string, value any, params []string, ctx RuleContext) bool {
	pv := ctx.GetPresenceVerifier()

	if pv == nil {
		return true
	}

	if len(params) < 2 {
		return true
	}

	table := params[0]
	column := params[1]

	var excludeID *string

	var idColumn *string

	if len(params) >= 3 && params[2] != "" && params[2] != "NULL" {
		s := params[2]
		excludeID = &s
	}

	if len(params) >= 4 && params[3] != "" {
		s := params[3]
		idColumn = &s
	}

	start := 4

	if excludeID == nil {
		start = 2
	}

	extras := parseExtras(params[start:])

	count := pv.GetCount(table, column, stringify(value), excludeID, idColumn, extras)

	return count == 0
}

func parseExtras(params []string) map[string]any {
	if len(params) == 0 {
		return nil
	}

	out := make(map[string]any)

	for i := 0; i+1 < len(params); i += 2 {
		out[params[i]] = params[i+1]
	}

	return out
}
