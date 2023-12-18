// Package lodash
// CodeAnalyser for Lodash GitHub repository for javascript/typescript files
package lodash

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v57/github"
	"github.com/kotvrt/files-letter-analyzer/alphabet"
	"log"
	"os"
	"strconv"
	"time"
)

var ErrRateLimited = errors.New("rate limited by GitHub")

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
	cfg               *Config
	executionDuration time.Duration
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
		log.Fatalf("fatal: github token hasn't be set - hint:`export GITHUB_TOKEN=<your-github-token>`")
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
	// Golang provides a context
	// this is a initialisation of a simple, empty context
	ctx := context.Background()

	metrics := make(map[string]int, len(alphabet.English))

	for index := 0; index < len(alphabet.English); index++ {
		letter := alphabet.English[index]
		//TODO: Would be better to have the caller have a handle over this via config
		time.Sleep(2 * time.Second)
		query := a.createSearchQueryForLetter(letter)

		// NOTE: GitHub Search API has a rate limit of up to:
		// - 30 reqs/min for authenticated users
		// - 10 reqs/min for unauthenticated users
		// Check out GitHub's Search Code doc for more info on limitations:
		// https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-code
		codeSearchResult, githubResponse, err := client.Search.Code(
			ctx,
			query,
			&github.SearchOptions{
				TextMatch: true,
			})

		// Means request has been rate limited by GitHub
		// Read more on GitHub API's rate limiting policy by following link below:
		// https://docs.github.com/en/rest/using-the-rest-api/troubleshooting-the-rest-api?apiVersion=2022-11-2
		if githubResponse.Response.StatusCode == 403 {
			rateLimitDuration := maybeFetchRateLimitDurationFromHeader(githubResponse)
			//TODO: upper limit to waiting time is a good candidate for a config knob
			if rateLimitDuration == nil || *rateLimitDuration > time.Minute*2 {
				// return special error to tell caller that metrics map can and should still be parsed
				// this way we guard the metrics that may have already been fetched prior to rate-limiting
				return ErrRateLimited, metrics
			}
			// wait out duration of rate limit and continue
			time.Sleep(*rateLimitDuration)
			// retry for the previous letter that was blocked by rate-limiting
			index = index - 1
			continue
		}

		if err != nil {
			return fmt.Errorf("error doing Github search: %w", err), nil
		}

		metrics[letter] = codeSearchResult.GetTotal()
	}

	return nil, metrics
}

func maybeFetchRateLimitDurationFromHeader(githubResponse *github.Response) *time.Duration {
	if githubResponse == nil {
		return nil
	}
	// When GitHub API's rate-limits the caller it will set the X-Ratelimit-Reset header's value
	// the value represents the time when the rate limit would be lifted
	// it's an integer representing Unix time since Epoch in seconds
	rateLimitExpiryInSecsSinceEpoch, err := strconv.Atoi(githubResponse.Header.Get("X-Ratelimit-Reset"))
	if err != nil {
		return nil
	}
	durationUntilExpiry := time.Until(time.Unix(int64(rateLimitExpiryInSecsSinceEpoch), 0))
	return &durationUntilExpiry
}

func (a CodeAnalyser) createSearchQueryForLetter(letter string) string {
	// The library is going to encode the search query as HTTP query parameter
	return fmt.Sprintf("%s language:JavaScript repo:%s",
		letter, a.cfg.GithubRepository)
}
