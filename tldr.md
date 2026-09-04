TLGC(1)

# NAME

TLGC - a tldr client written in Go

# SYNOPSIS

_tldr_ [OPTIONS] [PAGE]...

# DESCRIPTION

_tlgc_ is a command-line client for tldr-pages,
a community-driven collection of simplified man pages.
It downloads, caches, searches, renders, lints, and formats tldr pages.

_tlgc_ implements client specification v2.3.

# Options

Exactly one of the following operations may be active at a time.
When no operation is specified and a positional argument is given,
the argument is treated as a page name to look up.

## Operations

_tldr_ [PAGE]...
Look up and render the given tldr page.
If the page is not in the local cache,
it is fetched automatically (unless _-o_ is set).

_-r, --render_ <FILE>
Render a local tldr page file.
The file is parsed and displayed using the same renderer as a cached page.

_-b, --browse_ [PAGE]...
Open the page in the default web browser
instead of rendering it in the terminal.

_-u, --update_
Download the latest tldr-page archives
and refresh the local cache.

_-i, --info_
Display cache information: directory, size, age, and language statistics.

_--clean-cache_
Interactively delete the contents of the cache directory.
Prompts for confirmation before removing files.

_-s, --search_ <KEYWORD>
Search all cached pages (all languages) for pages containing the given keyword.
Results are displayed in a highlighted table.

_-l, --list_
List all pages available for the current platform (English only).

_-a, --list-all_
List all pages across all platforms (English only).

_--list-platforms_
List available platforms: common, linux, osx, windows, android.

_--list-languages_
List installed languages (i.e. languages that have a local cache).

_--lint_ <FILE|DIR>
Validate tldr pages against the tldr-pages style guide. Prints
each violation with its error code and line number.

_--format_ <FILE|DIR>
Reformat tldr pages to canonical style. Output goes to stdout
by default; use _--output_ or _--in-place_ to redirect it.

_--gen-config_
Print the default configuration as TOML to stdout.

_--config-path_
Print the default configuration file path to stdout.

_-v, --version_
Print the version string.

_-h, --help_
Print help information.

## Modifiers

Flags that modify the behavior of the commands above.
Many options are _modifiers_ that only apply to specific parent operations.

_-p, --platform_ <PLATFORM>
Override the auto-detected platform.
Valid values: linux, osx (alias: macos), windows, android, common.
Applies to: page lookup, _--browse_, _--list_, _--search_.

_-L, --language_ <CODE>
Override the auto-detected language(s).
May be specified multiple times.
Applies to: page lookup, _--browse_, _--search_, _--update_.

_-o, --offline_
Do not update the cache, even if it is stale.
Applies to: page lookup, _--browse_.

_-c, --compact_
Strip blank separator lines from rendered output.
Applies to: page lookup, _--render_.

_--no-compact_
Do not strip blank lines, Overrides _--compact_.
Applies to: page lookup, _--render_.

_-R, --raw_
Print the raw markdown source instead of rendering.
Applies to: page lookup, _--render_.

_--no-raw_
Render the page instead of printing raw markdown, Overrides _--raw_.
Applies to: page lookup, _--render_.

_--short-options_
Display short flags where possible (e.g. _-s_ instead of _--search_).
Applies to: page lookup, _--render_.

_--long-options_
Display long flags where possible (e.g. _--search_ instead of _-s_).
Applies to: page lookup, _--render_. This is the default.

_--edit_
Display a link to edit the page on GitHub.
Applies to: page lookup, _--render_.

_--tabular_
Display lint errors in a tabular format.
Applies to: _--lint_, _--format_.

_--ignore_ <CODES>
Comma-separated list of error codes to suppress (e.g. _TLDR005,TLDR019_).
Applies to: _--lint_, _--format_.

_--output_ <FILE>
Write formatted output to a file instead of stdout.
Applies to: _--format_.

_--in-place_
Rewrite files in place.
Applies to: _--format_.

_-q, --quiet_
Suppress all status messages, warnings, and informational output.

_--verbose_
Increase verbosity. May be specified up to twice for more detail.

_--color_ <WHEN>
Control when color is enabled.
Possible values: _auto_ (default), _always_, _never_.
When _auto_, color is enabled only when stdout is a terminal and the terminal supports it.

_--config_ <FILE>
Use an alternative configuration file instead of the default path.

# CONFIGURATION

tlgc uses a TOML configuration file at a platform-dependent location:

- Linux: ~/.config/tlgc/config.toml
- macOS: ~/Library/Application Support/tlgc/config.toml
- Windows: %AppData%/tlgc/config.toml

The path can be overridden with the _TLGC_CONFIG_ environment variable.
If no config file exists, built-in defaults are used silently.

Run _tldr --gen-config_ to print the default configuration.
See the full configuration reference in the docs/ directory.

# ENVIRONMENT

_TLGC_CONFIG_
Override the default configuration file path.

_LANG_, _LANGUAGE_
Used for automatic language detection when cache.languages is empty.
LANGUAGE is split on colons; for ll_CC locales,
both the full locale and the bare language code are tried.

_NO_COLOR_
When set (any value), disables colored output.
Equivalent to _--color never_.

# EXIT STATUS

0
Success.

1
An error occurred (invalid arguments, cache failure, config error, page not found, etc.).

# EXAMPLES

Look up the `tar` page:

	tldr tar

Force the Linux platform:

	tldr -p linux tar

Search for pages that contains the string `ngin`:

	tldr -s ngin

Lint a tldr page:

	tldr --lint path/to/page.md

Reformat a page in place:

	tldr --format --in-place pages/to/page.md

Open a page in the browser:

	tldr -b git

# SEE ALSO

tldr client specification
https://github.com/tldr-pages/tldr/blob/main/CLIENT-SPECIFICATION.md

tlgc's source code (report issues with the client here)
https://github.com/TheRootDaemon/tlgc

tldr-pages upstream (report issues with the pages here)
https://github.com/tldr-pages/tldr
