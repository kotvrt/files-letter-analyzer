package main

import (
	"fmt"
	"github.com/kotvrt/files-letter-analyzer/analyzer"
	"github.com/kotvrt/files-letter-analyzer/analyzer/lodash"
	"log"
)

func main() {
	var githubCodeAnalyzer analyzer.Analyser
	githubCodeAnalyzer = lodash.NewCodeAnalyser()

	err, _ := githubCodeAnalyzer.Analyse()
	if err != nil {
		log.Fatalln(fmt.Errorf("fatal error was encountered while trying to run the GitHub code analysis: %w", err))
	}

	//TODO: print metrics
	//var metricPrinter metrics.Printer
}
