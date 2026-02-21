package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// Promise represents the outcome of an agent's task as stated in the response.
type Promise string

const (
	PromiseComplete Promise = "COMPLETE"
	PromiseFailed   Promise = "FAILED"
	PromiseUnknown  Promise = "UNKNOWN"
)

var promiseRegex = regexp.MustCompile(`(?si)<promise>(.*?)</promise>`)

// ExtractPromise identifies the promise made by the agent in the response.
// It uses MarkdownSanitizer to ignore tags within code blocks.
func ExtractPromise(response string) (Promise, error) {
	s := NewMarkdownSanitizer()
	sanitized := s.StripCodeBlocks(response)

	matches := promiseRegex.FindAllStringSubmatch(sanitized, -1)
	if len(matches) == 0 {
		return PromiseUnknown, nil
	}

	// Use the first promise if multiple are present
	rawPromise := strings.TrimSpace(matches[0][1])
	upperPromise := strings.ToUpper(rawPromise)

	switch upperPromise {
	case "COMPLETE":
		return PromiseComplete, nil
	case "FAILED":
		return PromiseFailed, nil
	default:
		return PromiseUnknown, fmt.Errorf("invalid promise value: %q", rawPromise)
	}
}
