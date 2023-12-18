package lodash

import (
	"github.com/kotvrt/files-letter-analyzer/analyzer"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func Test_CallToAnalyserExecutesSuccessfully(t *testing.T) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	assert.Check(t, githubToken != "", "expecting github token not to be empty")
	assert.NilError(t, os.Unsetenv("GITHUB_TOKEN"))

	var githubCodeAnalyzer analyzer.Analyser
	githubCodeAnalyzer = NewCodeAnalyser(WithGithubToken(githubToken))

	err, metrics := githubCodeAnalyzer.Analyse()
	assert.NilError(t, err)
	assert.Check(t, len(metrics) > 0)
	assert.Check(t, metrics["A"] > 0)
}
