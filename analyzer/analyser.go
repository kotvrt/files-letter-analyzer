package analyser

// Analyser defines behaviour of classes implementing the interface.
// Contract is delivered through analyze() function.
// analyze() returns error if something goes wrong
// map
type Analyser interface {
	Analyse() (error, map[string]int)
}
