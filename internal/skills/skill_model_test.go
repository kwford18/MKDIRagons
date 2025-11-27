package skills_test

import (
	"bytes"
	"github.com/kwford18/MKDIRagons/internal/skills"
	"io"
	"os"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/stretchr/testify/assert"
)

// captureOutput is a helper function to capture stdout for testing Print methods
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		out <- buf.String()
	}()

	f()
	_ = w.Close()
	os.Stdout = stdout
	return <-out
}

func TestSkill_GetEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		skillName string
		expected  string
	}{
		{"Normal Name", "Stealth", "skills/Stealth"},
		{"Empty Name", "", "skills/"},
		{"Lowercase Name", "arcana", "skills/arcana"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &skills.Skill{Name: tt.skillName}
			assert.Equal(t, tt.expected, s.GetEndpoint())
		})
	}
}

func TestSkill_Print(t *testing.T) {
	s := &skills.Skill{
		Name:       "Acrobatics",
		Bonus:      5,
		Ability:    core.Dexterity,
		Proficient: true,
		Expertise:  false,
	}

	output := captureOutput(func() {
		s.Print()
	})

	// Assert that the output contains the specific formatted lines
	assert.Contains(t, output, "- Name: Acrobatics")
	assert.Contains(t, output, "- Value: 5")
	assert.Contains(t, output, "- Proficient: true")
	assert.Contains(t, output, "- Expertise: false")

	// Check for Ability line using the String() method of the enum
	assert.Contains(t, output, "- Ability: Dexterity")
}

func TestSkillList_GetEndpoint(t *testing.T) {
	sl := &skills.SkillList{}
	assert.Equal(t, "skills/", sl.GetEndpoint())
}

func TestSkillList_Print(t *testing.T) {
	// Populate a few skills to verify the loop works correctly
	sl := &skills.SkillList{
		Athletics: skills.Skill{Name: "Athletics", Bonus: 4},
		Stealth:   skills.Skill{Name: "Stealth", Bonus: 8},
		// Zero-value skills will have empty names and 0 bonus,
		// but checking a few key ones ensures the mapping is correct.
	}

	output := captureOutput(func() {
		sl.Print()
	})

	// Verify specific formatted lines exist
	// %-20s pads the name, so we check for the name followed by spaces and the value
	assert.Contains(t, output, "Athletics:           4")
	assert.Contains(t, output, "Stealth:             8")

	// Verify that at least one other skill (zero value) was printed
	assert.Contains(t, output, "0")

	// Ensure we have 18 lines (one for each skill in the struct)
	// We count newlines to verify all skills in the hardcoded list were iterated
	lineCount := 0
	for _, char := range output {
		if char == '\n' {
			lineCount++
		}
	}
	assert.Equal(t, 18, lineCount, "Should print exactly 18 lines for the 18 skills")
}
