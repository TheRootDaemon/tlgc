# tlgc

A tldr client written in Go.

## Usage

See `tldr --help` for all options.

## Configuration

tlgc can be customized with a [TOML](https://toml.io) configuration file. To get the default path for your system, run:

```sh
tldr --config-path
```

To generate a default config file, run:

```sh
tldr --gen-config > "$(tldr --config-path)"
```

or copy the below example.

### Configuration options

```toml
[cache]
# Override the cache directory ('~' will be expanded to your home directory).
dir = "~/.cache/tlgc"
# Override the base URL used for downloading tldr pages.
# The mirror must provide files with the same names as the official tldr pages repository:
# mirror/tldr.sha256sums            must point to the SHA256 checksums of all assets
# mirror/tldr-pages.LANGUAGE.zip    must point to a zip archive that contains platform directories with pages in LANGUAGE
mirror = "https://github.com/tldr-pages/tldr/releases/latest/download"
# Automatically update the cache if it's older than max_age hours.
auto_update = true
# Perform the automatic update after the page is shown (the default is to update first, then show the page).
defer_auto_update = false
max_age = 336 # 336 hours = 2 weeks
# Specify a list of desired page languages. If it's empty, languages specified in
# the LANG and LANGUAGE environment variables are downloaded.
# English is implied and will always be downloaded.
# You can see a list of language codes here: https://github.com/tldr-pages/tldr
# Example: ["de", "pl"]
languages = []

[output]
# Show the title in the rendered page.
show_title = true
# Show the platform name ('common', 'linux', etc.) in the title.
platform_title = false
# Prefix descriptions of examples with hyphens.
show_hyphens = false
# Display a link to edit the shown page on GitHub.
edit_link = false
# Use a custom string instead of a hyphen.
example_prefix = "- "
# Set the max line length. 0 means wrapping is disabled.
# If a line is longer than this value, it will be split into multiple lines.
line_length = 0
# Strip blank separator lines from output.
compact = false
# In option placeholders, show the specified option style.
# Example: {{[-s|--long]}}
# short  : -s
# long   : --long
# both   : [-s|--long]
option_style = "long"
# Print pages in raw markdown.
raw_markdown = false

# Number of spaces to put before each line of the page.
[indent]
# Command name.
title = 2
# Command description.
description = 2
# Descriptions of examples.
bullet = 2
# Example command invocations.
example = 4

# Style for the title of the page (command name).
[style.title]
# Fixed colors:       "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white", "default",
#                     "bright_black", "bright_red", "bright_green", "bright_yellow", "bright_blue",
#                     "bright_magenta", "bright_cyan", "bright_white"
# 256color ANSI code: "color256:50"
# RGB:                "rgb:0,255,255"
# Hex:                "#ffffff"
color = "magenta"
background = "default"
bold = true
underline = false
italic = false
dim = false
strikethrough = false

# Style for the description of the page.
[style.description]
color = "magenta"
background = "default"
bold = false
underline = false
italic = false
dim = false
strikethrough = false

# Style for descriptions of examples.
[style.bullet]
color = "green"
background = "default"
bold = false
underline = false
italic = false
dim = false
strikethrough = false

# Style for command examples.
[style.example]
color = "cyan"
background = "default"
bold = false
underline = false
italic = false
dim = false
strikethrough = false

# Style for URLs inside the description.
[style.url]
color = "red"
background = "default"
bold = false
underline = false
italic = true
dim = false
strikethrough = false

# Style for text surrounded by backticks (`).
[style.inline_code]
color = "yellow"
background = "default"
bold = false
underline = false
italic = true
dim = false
strikethrough = false

# Style for placeholders inside command examples.
[style.placeholder]
color = "red"
background = "default"
bold = false
underline = false
italic = true
dim = false
strikethrough = false
```

## Linting

tlgc implements a built-in linter and a formatter for tldr pages.
It is intended to be a successor of [tldr-lint](https://github.com/tldr-pages/tldr-lint).

Validate one or more pages or directories:

```sh
tldr --lint path/to/page.md
```

Reformat pages to canonical style (stdout by default):

```sh
tldr --format path/to/page.md
```

Write the result back to the original file or to a new file:

```sh
tldr --format --in-place path/to/page.md
tldr --format --output path/to/result.md path/to/page.md
```

Modifiers:

- `--tabular` — display errors in a tabular format
- `--ignore TLDR005,TLDR019` — suppress specific error codes

### Error codes

#### Content rules

| Code    | Description                                                                           |
| ------- | ------------------------------------------------------------------------------------- |
| TLDR001 | File should contain no leading whitespace.                                            |
| TLDR002 | A single space should precede a sentence.                                             |
| TLDR003 | Descriptions should start with a capital letter.                                      |
| TLDR004 | Command descriptions should end in a period.                                          |
| TLDR005 | Example descriptions should end in a colon with no trailing characters.               |
| TLDR006 | Command name and description should be separated by an empty line.                    |
| TLDR007 | Example descriptions should be surrounded by empty lines.                             |
| TLDR008 | File should contain no trailing whitespace.                                           |
| TLDR009 | Page should contain a newline at end of file.                                         |
| TLDR010 | Only Unix-style line endings allowed.                                                 |
| TLDR011 | Page should never contain more than a single empty line.                              |
| TLDR012 | Page should contain no tabs.                                                          |
| TLDR013 | Title should be alphanumeric with dashes, underscores, spaces, or allowed characters. |
| TLDR014 | Page should contain no trailing whitespace.                                           |
| TLDR015 | Example descriptions should start with a capital letter.                              |
| TLDR016 | Label for information link should be spelled exactly `More information: `.            |
| TLDR017 | Information link should be surrounded with angle brackets.                            |
| TLDR018 | Page should only include a single information link.                                   |
| TLDR019 | Page should only include a maximum of 8 examples.                                     |
| TLDR020 | Label for additional notes should be spelled exactly `Note: `.                        |
| TLDR021 | Command example should not begin or end in whitespace.                                |

#### Hint rules

| Code    | Description                                                                              |
| ------- | ---------------------------------------------------------------------------------------- |
| TLDR101 | Command description probably not properly annotated.                                     |
| TLDR102 | Example description probably not properly annotated.                                     |
| TLDR103 | Command example is missing its closing backtick.                                         |
| TLDR104 | Example descriptions should prefer infinitive tense (e.g. write) over present or gerund. |
| TLDR105 | There should be only one command per example.                                            |

#### Filename rules

| Code    | Description                                                                                  |
| ------- | -------------------------------------------------------------------------------------------- |
| TLDR106 | Page title should start with a hash (`#`).                                                   |
| TLDR107 | File name should end with `.md` extension.                                                   |
| TLDR108 | File name should not contain whitespace.                                                     |
| TLDR109 | File name should be lowercase.                                                               |
| TLDR110 | Command example should not be empty.                                                         |
| TLDR111 | File name should not contain any Windows-forbidden character.                                |
| TLDR112 | Terms `stdin`, `stdout`, `stderr`, and `regex` should be lowercase and wrapped in backticks. |
