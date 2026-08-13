package app

import (
	"io/fs"
	"os"
	"path/filepath"
)

// collectFiles expands every input path into a flat, ordered list of page
// files.
//
// Each path is resolved with collectPathFiles:
// an individual file is taken as-is regardless of its extension
// (a non-.md file passed directly is still linted and reported as TLDR107),
// while a directory is walked recursively
// to gather every .md file beneath it.
// Files are returned in the order the paths were given,
// with directory contents ordered as filepath.WalkDir yields them.
//
// It returns an error and no files if any input path cannot be stat'd.
func collectFiles(paths []string) ([]string, error) {
	var files []string

	for _, path := range paths {
		pathFiles, err := collectPathFiles(path)
		if err != nil {
			return nil, err
		}
		files = append(files, pathFiles...)
	}

	return files, nil
}

// collectPathFiles expands a single path into the list of page files it designates.
//
// A path naming an existing file is returned as-is,
// whatever its extension, so that callers can lint explicitly named non-.md files.
// A path naming a directory is walked recursively,
// collecting only entries whose extension is ".md";
// subdirectories are descended automatically,
// and files of any other type are skipped.
// An empty directory yields an empty list.
//
// It returns an error if the path does not exist or cannot be stat'd,
// or if the directory walk fails partway through.
func collectPathFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	if err = filepath.WalkDir(
		path,
		func(entry string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(entry) == ".md" {
				files = append(files, entry)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return files, nil
}
