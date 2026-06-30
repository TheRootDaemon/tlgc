// Package render generates terminal-formatted output for tldr pages.
//
// It parses markdown-formatted tldr pages and renders them as
// colored, wrapped terminal text with syntax highlighting for command examples.
//
// Text is automatically wrapped to fit within a configurable maximum line width.
// Color output can be disabled with the WithColor(false) option.
package render
