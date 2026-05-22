package routing

// trait of the same name.
//
// In the upstream framework this is a trait used by [Route] and [RouteRegistrar] to expose
// the WhereAlpha / WhereNumber / WhereUuid family of helpers. In Go traits are
// realized as embedded structs: any type that embeds this struct gains the
// helpers, provided it also satisfies [whereTarget] (i.e. has a `Where` method
// the helpers can delegate to).
type CreatesRegularExpressionRouteConstraints struct {
	target whereTarget
}

// whereTarget is the contract a host type must satisfy to participate in the
// helpers below. Host types call Bind once during construction so the helpers
// can reach back into Where. The return value is ignored by the helpers but
// preserved so Route.Where can stay chainable.
type whereTarget interface {
	Where(name, expression string) *Route
}

// Bind wires the embedding host type into the helper. Hosts call this from
// their constructor, e.g. r.constraints.Bind(r).
func (c *CreatesRegularExpressionRouteConstraints) Bind(t whereTarget) { c.target = t }

// WhereAlpha constrains the named parameters to alphabetic characters only.
func (c *CreatesRegularExpressionRouteConstraints) WhereAlpha(parameters ...string) {
	c.assignExpression(parameters, `[a-zA-Z]+`)
}

// WhereAlphaNumeric constrains the named parameters to alphanumeric characters.
func (c *CreatesRegularExpressionRouteConstraints) WhereAlphaNumeric(parameters ...string) {
	c.assignExpression(parameters, `[a-zA-Z0-9]+`)
}

// WhereNumber constrains the named parameters to digits.
func (c *CreatesRegularExpressionRouteConstraints) WhereNumber(parameters ...string) {
	c.assignExpression(parameters, `[0-9]+`)
}

// WhereUlid constrains the named parameters to ULIDs (Crockford's base32).
func (c *CreatesRegularExpressionRouteConstraints) WhereUlid(parameters ...string) {
	c.assignExpression(parameters, `[0-7][0-9a-hjkmnp-tv-zA-HJKMNP-TV-Z]{25}`)
}

// WhereUuid constrains the named parameters to RFC 4122 UUIDs.
func (c *CreatesRegularExpressionRouteConstraints) WhereUuid(parameters ...string) {
	c.assignExpression(parameters, `[\da-fA-F]{8}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{12}`)
}

// WhereIn constrains the named parameters to one of the supplied literal
// values. Values are joined with "|" inside an alternation.
func (c *CreatesRegularExpressionRouteConstraints) WhereIn(parameter string, values []string) {
	expr := ""

	for i, v := range values {
		if i > 0 {
			expr += "|"
		}

		expr += v
	}

	c.target.Where(parameter, expr)
}

func (c *CreatesRegularExpressionRouteConstraints) assignExpression(parameters []string, expression string) {
	for _, p := range parameters {
		c.target.Where(p, expression)
	}
}
