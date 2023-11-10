package jsoncall

// TraceVariables provides the names and references to bool variables
// that control debugging output.
func TraceVariables() map[string]*bool {
	return map[string]*bool{
		"data":  &TraceData,
		"calls": &TraceCalls,
		"init":  &TraceInit,
	}
}
