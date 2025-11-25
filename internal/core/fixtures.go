package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadFixture loads a JSON fixture file from testdata directory
// and returns it as a map for use in mock HTTP responses
func LoadFixture(t *testing.T, filename string) map[string]interface{} {
	t.Helper()

	path := filepath.Join("testdata/fixtures", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read fixture file: %s", path)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Failed to unmarshal fixture: %s", path)

	return result
}

// LoadFixtureRaw loads raw JSON bytes from testdata directory
// Useful for testing malformed JSON or custom unmarshaling
func LoadFixtureRaw(t *testing.T, filename string) []byte {
	t.Helper()

	path := filepath.Join("testdata/fixtures", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read fixture file: %s", path)

	return data
}

// LoadFixtureInto loads a JSON fixture and unmarshals it directly into the provided struct
// This is useful for testing the full unmarshaling pipeline
func LoadFixtureInto(t *testing.T, filename string, target interface{}) {
	t.Helper()

	data := LoadFixtureRaw(t, filename)
	err := json.Unmarshal(data, target)
	require.NoError(t, err, "Failed to unmarshal fixture into target: %s", filename)
}
