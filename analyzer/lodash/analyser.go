// Package lodash
// CodeAnalyser for Lodash GitHub repository for javascript/typescript files
package lodash

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v57/github"
	"log"
	"os"
)

type Config struct {
	GithubToken         string
	GithubRepository    string
	GithubBaseUrl       string
	GithubFileExtension string
}

// CodeAnalyser implements the Analyser interface
// The implementing method 'Analyse()' searches for occurrences of letters in the content
// of .js/.ts files in Lodash GitHub repository: https://github.com/lodash/lodash
type CodeAnalyser struct {
	cfg *Config
}

type AnalyzerOption func(*Config)

func WithGithubToken(token string) AnalyzerOption {
	return func(a *Config) {
		a.GithubToken = token
	}
}

func WithGithubRepository(token string) AnalyzerOption {
	return func(a *Config) {
		a.GithubToken = token
	}
}

// NewCodeAnalyser accepts an arbitrary number of configuration options that would be applied to
// the returned instance. You can use these to override default configuration
// from inside the code, e.g. setting up mock configuration for tests.
func NewCodeAnalyser(initOptions ...AnalyzerOption) CodeAnalyser {
	cfg := configFromEnvironment()

	// options will override the environment variables
	for _, option := range initOptions {
		option(cfg)
	}

	if cfg.GithubToken == "" {
		log.Println("warning: github token hasn't be set; this will result in longer processing times")
	}

	return CodeAnalyser{
		cfg: cfg,
	}
}

func configFromEnvironment() *Config {
	return &Config{
		GithubToken:      os.Getenv("GITHUB_TOKEN"),
		GithubRepository: "lodash/lodash",
		GithubBaseUrl:    "https://api.github.com",
	}
}

// Analyse implicitly sets the contract between CodeAnalyser struct
// and the Analyzer interface
func (a CodeAnalyser) Analyse() (error, map[string]int) {
	client := github.NewClient(nil).WithAuthToken(a.cfg.GithubToken)
	if client == nil {
		return errors.New("fatal: something went wrong with creating github client"), nil
	}
	ctx := context.Background()
	query := a.createSearchQueryForLetter("a")

	// NOTE: GitHub Search API has a rate limit of up to:
	// - 30 reqs/min for authenticated users
	// - 10 reqs/min for unauthenticated users
	// Checkout GitHub's Search Code doc for more info on limitations:
	// https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-code
	_, _, err := client.Search.Code(
		ctx,
		query,
		&github.SearchOptions{
			// Order: 'desc' by default
			TextMatch: true,
		})
	//TODO: Handle rate limiting:
	// https://docs.github.com/en/rest/using-the-rest-api/troubleshooting-the-rest-api?apiVersion=2022-11-28
	if err != nil {
		return fmt.Errorf("error making Github search: %w", err), nil
	}
	// process the letter occurrences
	// return the results
	return nil, nil
}

func (a CodeAnalyser) createSearchQueryForLetter(letter string) string {
	return fmt.Sprintf("%s+language:JavaScript+AND+language:TypeScript+repo:%s",
		letter, a.cfg.GithubRepository)
}
