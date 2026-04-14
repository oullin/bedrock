package rules

// registerAll is called once from the registry init() function to register
// every built-in rule.
func registerAll() {
	init_accepted()
	init_required()
	init_filled()
	init_present()
	init_missing()
	init_required_markers()
	init_prohibited()
	init_excluded()
	init_type()
	init_string()
	init_regex()
	init_email()
	init_size()
	init_numeric()
	init_comparison()
	init_date()
	init_array()
	init_file()
	init_password()
	init_database()
}
