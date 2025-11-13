# Testing with Fixtures

This document explains how to use test fixtures in the MKDIRagons project.

## Setup

### 1. Create the fixtures

```bash
chmod +x setup_fixtures.sh
./setup_fixtures.sh internal/core/
```

This creates `internal/core/testdata/` with real API responses.

### 2. Verify fixtures were created

```bash
ls -la internal/core/testdata/
```

You should see:
- `wizard.json`
- `dwarf.json`
- `fireball.json`
- `padded-armor.json`

## Running Tests

### Run all tests
```bash
go test -v ./internal/core
```

### Run only fixture-based tests
```bash
go test -v ./internal/core -run RealJSON
```

### Run specific fixture test
```bash
go test -v ./internal/core -run TestFetchJSON_RealJSON_Wizard
```

### Run with coverage
```bash
go test -cover ./internal/core
go test -coverprofile=coverage.out ./internal/core
go tool cover -html=coverage.out
```

## Test Organization

### Files

- `fetch_json.go` - Production code with refactored `FetchJSON`
- `fetch_json_test.go` - Unit tests with inline mocks
- `fetch_json_fixtures_test.go` - Integration tests using real JSON
- `fixtures.go` - Helper functions for loading fixtures
- `testdata/` - Real API response JSON files

### Test Types

**Unit Tests (inline mocks):**
- Fast and focused
- Test specific behaviors
- Easy to modify for edge cases
- Located in `fetch_json_test.go`

**Integration Tests (real JSON):**
- Realistic and comprehensive
- Test full unmarshaling pipeline
- Catch struct field mismatches
- Located in `fetch_json_fixtures_test.go`

## Adding New Fixtures

### 1. Get the real API response

```bash
curl https://www.dnd5eapi.co/api/2014/classes/fighter | jq > internal/core/testdata/fighter.json
```

### 2. Create a test

```go
func TestFetchJSON_RealJSON_Fighter(t *testing.T) {
    mockResponse := loadFixture(t, "fighter.json")
    
    server := createMockServer(t, http.StatusOK, mockResponse)
    defer server.Close()

    class := &Class{}
    err := fetchJSONWithClient(http.DefaultClient, server.URL+"/", class, "fighter")
    require.NoError(t, err)
    
    assert.Equal(t, "fighter", class.Index)
    assert.Equal(t, 10, class.HitDie)
}
```

### 3. Add to table-driven test

Add your new test case to `TestFetchJSON_AllRealJSON` in `fetch_json_fixtures_test.go`.

## Helper Functions

### `LoadFixture(t, filename)`
Loads a JSON file and returns it as `map[string]interface{}`

```go
mockResponse := loadFixture(t, "wizard.json")
```

### `LoadFixtureRaw(t, filename)`
Loads raw JSON bytes (useful for testing malformed JSON)

```go
rawJSON := loadFixtureRaw(t, "wizard.json")
```

### `LoadFixtureInto(t, filename, target)`
Loads and unmarshals directly into a struct

```go
class := &Class{}
loadFixtureInto(t, "wizard.json", class)
```

## Best Practices

### ✅ DO:
- Use real JSON for complex nested structures
- Use inline mocks for simple unit tests
- Test error cases with inline mocks
- Keep fixtures up-to-date with API changes
- Validate all important fields in fixture tests

### ❌ DON'T:
- Modify fixture files manually (regenerate from API)
- Use fixtures for testing error responses
- Create fixtures for every possible input
- Commit large (>100KB) fixture files without review

## Troubleshooting

### Fixture file not found
```
Error: Failed to read fixture file: testdata/wizard.json
```

**Solution:** Run `./setup_fixtures.sh` to create fixtures.

### JSON unmarshal error
```
Error: Failed to unmarshal fixture: invalid character...
```

**Solution:** Validate JSON with `jq`:
```bash
cat internal/core/testdata/wizard.json | jq
```

### Test failing with real JSON but passing with inline mocks

This likely means your struct tags don't match the real API response. Check:
1. Field names match JSON keys (case-sensitive)
2. JSON tags are correct (`json:"field_name"`)
3. All required fields are present in struct

## Example Test Run

```bash
$ go test -v ./internal/core -run RealJSON

=== RUN   TestFetchJSON_RealJSON_Wizard
--- PASS: TestFetchJSON_RealJSON_Wizard (0.01s)
=== RUN   TestFetchJSON_RealJSON_Dwarf
--- PASS: TestFetchJSON_RealJSON_Dwarf (0.01s)
=== RUN   TestFetchJSON_RealJSON_Fireball
--- PASS: TestFetchJSON_RealJSON_Fireball (0.00s)
=== RUN   TestFetchJSON_RealJSON_PaddedArmor
--- PASS: TestFetchJSON_RealJSON_PaddedArmor (0.00s)
PASS
ok      github.com/kwford18/MKDIRagons/internal/core    0.023s
```

## Updating Fixtures

When the D&D API changes:

```bash
# Update all fixtures
./setup_fixtures.sh

# Or update individual fixtures
curl https://www.dnd5eapi.co/api/2014/classes/wizard | jq > internal/core/testdata/wizard.json
```

Then run tests to ensure structs still match:
```bash
go test -v ./internal/core -run RealJSON
```