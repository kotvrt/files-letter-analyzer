package main

import (
	"errors"
	"fmt"
	"github.com/kotvrt/files-letter-analyzer/analyzer"
	"github.com/kotvrt/files-letter-analyzer/analyzer/lodash"
	maps "github.com/tg/gosortmap"
	"log"
	"os"
	"time"
)

func main() {
	var githubCodeAnalyzer analyzer.Analyser
	githubCodeAnalyzer = lodash.NewCodeAnalyser()

	fmt.Println("Fetching alphabet metrics...")

	analysisBeginning := time.Now()
	err, metrics := githubCodeAnalyzer.Analyse()
	sortedMetrics := maps.ByValueDesc(metrics)

	if err != nil {
		if errors.Is(err, lodash.ErrRateLimited) {
			log.Printf("incomplete results, one or multiple calls to GitHub API have been rate limited")
			printMetricsElapsedTimeAndExit(sortedMetrics, analysisBeginning)
		}
		log.Fatalln(fmt.Errorf("fatal error was encountered "+
			"while trying to run the GitHub code analysis: %w", err))
	}
	printMetricsElapsedTimeAndExit(sortedMetrics, analysisBeginning)
}

func printMetricsElapsedTimeAndExit(metrics maps.Items, analysisBeginning time.Time) {
	for _, metric := range metrics {
		log.Printf("letter: %s occurrs in %d files\n---", metric.Key, metric.Value)
	}
	log.Printf("elapsed: %v\n", time.Since(analysisBeginning))
	os.Exit(0)
}
