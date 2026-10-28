package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runGrep(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmdArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func makeInputFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.txt")
	assert.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestGrepMatchesPattern(t *testing.T) {
	input := makeInputFile(t, "first\nauth event\nlast\n")

	output, err := runGrep(t, "auth", input)
	assert.NoError(t, err)
	assert.Equal(t, "auth event\n", output)
}

func TestGrepPrintsMatchingLineNumbers(t *testing.T) {
	input := makeInputFile(t, "first\nauth event\nlast\nauth failed\n")

	output, err := runGrep(t, "-n", "auth", input)
	assert.NoError(t, err)
	assert.Equal(t, "2. auth event\n4. auth failed\n", output)
}

func TestGrepFixedStringMatchesSubstring(t *testing.T) {
	input := makeInputFile(t, "prefix auth suffix\nauthentic\nother\n")

	output, err := runGrep(t, "-F", "auth", input)
	assert.NoError(t, err)
	assert.Equal(t, "prefix auth suffix\nauthentic\n", output)
}

func TestGrepCountsNonMatchingLines(t *testing.T) {
	input := makeInputFile(t, "auth event\nother\nauth failed\n")

	output, err := runGrep(t, "-c", "-v", "auth", input)
	assert.NoError(t, err)
	assert.Equal(t, "1\n", output)
}
