package md_note

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/shurcooL/github_flavored_markdown"
)

func CompileMD(input string) string {
	unsafeHTML := github_flavored_markdown.Markdown([]byte(input))
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	html := policy.SanitizeBytes(unsafeHTML)
	return string(html[:])
}
