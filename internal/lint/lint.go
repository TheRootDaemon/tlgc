package lint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Error represents a single lint violation.
type Error struct {
	Code        string // e.g. "TLDR001"
	Line        int    // 1-indexed line number; 0 for file-level errors
	Description string
}

// Result holds all lint violations found.
type Result struct {
	Errors []Error
}

// ErrorCodes maps code to human-readable description.
var ErrorCodes = map[string]string{
	"TLDR001": "File should contain no leading whitespace",
	"TLDR002": "A single space should precede a sentence",
	"TLDR003": "Descriptions should start with a capital letter",
	"TLDR004": "Command descriptions should end in a period",
	"TLDR005": "Example descriptions should end in a colon with no trailing characters",
	"TLDR006": "Command name and description should be separated by an empty line",
	"TLDR007": "Example descriptions should be surrounded by empty lines",
	"TLDR008": "File should contain no trailing whitespace",
	"TLDR009": "Page should contain a newline at end of file",
	"TLDR010": "Only Unix-style line endings allowed",
	"TLDR011": "Page never contains more than a single empty line",
	"TLDR012": "Page should contain no tabs",
	"TLDR013": "Title should be alphanumeric with dashes, underscores, spaces or allowed characters",
	"TLDR014": "Page should contain no trailing whitespace",
	"TLDR015": "Example descriptions should start with a capital letter",
	"TLDR016": "Label for information link should be spelled exactly `More information: `",
	"TLDR017": "Information link should be surrounded with angle brackets",
	"TLDR018": "Page should only include a single information link",
	"TLDR019": "Page should only include a maximum of 8 examples",
	"TLDR020": "Label for additional notes should be spelled exactly `Note: `",
	"TLDR021": "Command example should not begin or end in whitespace",
	"TLDR101": "Command description probably not properly annotated",
	"TLDR102": "Example description probably not properly annotated",
	"TLDR103": "Command example is missing its closing backtick",
	"TLDR104": "Example descriptions should prefer infinitive tense (e.g. write) over present (e.g. writes) or gerund (e.g. writing)",
	"TLDR105": "There should be only one command per example",
	"TLDR106": "Page title should start with a hash ('#')",
	"TLDR107": "File name should end with .md extension",
	"TLDR108": "File name should not contain whitespace",
	"TLDR109": "File name should be lowercase",
	"TLDR110": "Command example should not be empty",
	"TLDR111": "File name should not contain any Windows-forbidden character",
	"TLDR112": "Terms `stdin`, `stdout`, `stderr`, and `regex` should be lowercase and wrapped in backticks",
}

// rule couples a TLDR code with the check that reports it.
type rule struct {
	code  string
	check func(*parsedPage, *Result)
}

// contentRules lists every rule that runs on page content, in ascending
// TLDR code order. Filename-only rules (TLDR107, TLDR108, TLDR109, TLDR111)
// are excluded because they require a filename.
var contentRules = []rule{
	{"TLDR001", checkLeadingWhitespace},
	{"TLDR002", checkSpaceAfterPrefix},
	{"TLDR003", checkDescriptionStartsWithCapital},
	{"TLDR004", checkDescriptionEndsWithPeriod},
	{"TLDR005", checkExampleDescriptionEndsWithColon},
	{"TLDR006", checkTitleDescriptionSeparator},
	{"TLDR007", checkExampleDescriptionSurroundedByBlankLines},
	{"TLDR008", checkNoTrailingWhitespaceAtEOF},
	{"TLDR009", checkEndsWithNewline},
	{"TLDR010", checkUnixLineEndings},
	{"TLDR011", checkConsecutiveBlankLines},
	{"TLDR012", checkNoTabs},
	{"TLDR013", checkTitleCharacters},
	{"TLDR014", checkTrailingWhitespace},
	{"TLDR015", checkExampleDescriptionStartsWithCapital},
	{"TLDR016", checkInformationLinkLabel},
	{"TLDR017", checkInformationLinkBrackets},
	{"TLDR018", checkSingleInformationLink},
	{"TLDR019", checkMaximumExampleCount},
	{"TLDR020", checkNoteLabelFormat},
	{"TLDR021", checkCommandWhitespace},
	{"TLDR101", checkCommandDescriptionAnnotated},
	{"TLDR102", checkExampleDescriptionAnnotated},
	{"TLDR103", checkCommandClosingBacktick},
	{"TLDR104", checkInfinitiveTense},
	{"TLDR105", checkSingleCommandPerExample},
	{"TLDR106", checkTitleHash},
	{"TLDR110", checkCommandNotEmpty},
	{"TLDR112", checkStandardTermsInBackticks},
}

// filenameRule couples a TLDR code with the check that reports it.
type filenameRule struct {
	code  string
	check func(string, *Result)
}

// filenameRules lists every rule that runs on a page's filename, in
// ascending TLDR code order.
var filenameRules = []filenameRule{
	{"TLDR107", checkFileExtension},
	{"TLDR108", checkFilenameWhitespace},
	{"TLDR109", checkFilenameLowercase},
	{"TLDR111", checkForbiddenFilenameCharacters},
}

// Lint runs all applicable rules on the page read from f
// and returns the violations found.
// Rules may inspect the file's content or name.
// Any rule code listed in ignore is skipped.
func Lint(
	f *os.File,
	ignore ...string,
) (*Result, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return &Result{}, err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return &Result{}, err
	}

	return lint(
		filepath.Base(f.Name()),
		content,
		ignore...,
	), nil
}

// lint runs all applicable rules on the page
// with the given filename
// and content and returns the violations found.
// Any rule code listed in ignore is skipped.
func lint(
	name string,
	content []byte,
	ignore ...string,
) *Result {
	r := &Result{}

	skipped := make(map[string]struct{}, len(ignore))
	for _, code := range ignore {
		skipped[code] = struct{}{}
	}

	if page := parse(string(content)); page != nil {
		for _, rl := range contentRules {
			if _, ok := skipped[rl.code]; ok {
				continue
			}
			rl.check(page, r)
		}
	}

	for _, rl := range filenameRules {
		if _, ok := skipped[rl.code]; ok {
			continue
		}
		rl.check(name, r)
	}

	return r
}

// addError is a convenience helper used by rules.
func addError(r *Result, code string, line int) {
	desc := ErrorCodes[code]
	if desc == "" {
		desc = code
	}
	r.Errors = append(
		r.Errors,
		Error{
			Code:        code,
			Line:        line,
			Description: desc,
		},
	)
}

// String returns the error as a formatted string
// containing its code, line number, and description.
func (e Error) String() string {
	return fmt.Sprintf(
		"%s:%d %s",
		e.Code,
		e.Line,
		e.Description,
	)
}
