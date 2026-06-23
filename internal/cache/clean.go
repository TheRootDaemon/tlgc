package cache

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
)

func (c *Cache) Clean() error {
	entries, err := getEntries(c.dir)
	if err != nil {
		return err
	}

	if entries == nil {
		logger.Info("cache does not exist, skipping...")
		return nil
	}

	var log strings.Builder
	log.WriteString("removing following files...\n")
	for _, entry := range entries {
		name := entry.Name()
		fmt.Fprintf(&log, "\n%q", name)
	}

	logger.InfoStart(
		"%s\nproceed with cleaning: [Y/n] ",
		log.String(),
	)

	cleanCache := parseInput(bufio.NewReader(os.Stdin))
	if !cleanCache {
		logger.InfoEnd("aborted")
		return nil
	}

	logger.Info("cleaning...")
	for _, entry := range entries {
		if err := os.RemoveAll(
			filepath.Join(
				c.dir,
				entry.Name(),
			),
		); err != nil {
			return fmt.Errorf(
				"remove %q: %w",
				entry.Name(),
				err,
			)
		}
	}

	logger.InfoEnd("done...")

	return nil
}

func getEntries(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	if len(entries) == 0 {
		return nil, nil
	}

	return entries, nil
}

func parseInput(reader *bufio.Reader) bool {
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return true
	}

	switch {
	case input == "yes":
		return true
	case input[0] == 'y':
		return true
	case input[0] == 'Y':
		return true
	default:
		return false
	}
}
