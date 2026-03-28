package skill

import (
	"runtime"
	"strings"
)

// Entry is a fully parsed skill ready for use.
type Entry struct {
	// Name is the canonical skill identifier (from frontmatter or directory name).
	Name string
	// Description is the human-readable summary from frontmatter.
	Description string
	// Content is the Markdown body (after frontmatter).
	Content string
	// Frontmatter holds all parsed frontmatter fields.
	Frontmatter Frontmatter
	// Source indicates where this skill was loaded from.
	Source string
	// FilePath is the absolute path to the SKILL.md file.
	FilePath string
}

// EligibilityContext carries runtime context for skill filtering.
type EligibilityContext struct {
	// WorkspaceSkillNames is the explicit list from session (empty = accept all).
	WorkspaceSkillNames []string
}

// IsEligible returns true if this skill should be injected into the system prompt.
// Mirrors openclaw TS filterSkillEntries logic.
func (e *Entry) IsEligible(ctx EligibilityContext) bool {
	// OS filter
	if len(e.Frontmatter.OS) > 0 && !matchesCurrentOS(e.Frontmatter.OS) {
		return false
	}

	// Binary requirements
	for _, bin := range e.Frontmatter.Requires {
		if !hasBinary(bin) {
			return false
		}
	}

	// If session explicitly lists skills, only include those (plus always=true)
	if len(ctx.WorkspaceSkillNames) > 0 && !e.Frontmatter.Always {
		found := false
		for _, n := range ctx.WorkspaceSkillNames {
			if n == e.Name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// ShouldInjectIntoPrompt returns true if this skill's content should appear in the system prompt.
func (e *Entry) ShouldInjectIntoPrompt() bool {
	return !e.Frontmatter.DisableModelInvocation
}

func matchesCurrentOS(osList []string) bool {
	current := runtime.GOOS // "darwin", "linux", "windows"
	for _, o := range osList {
		if strings.EqualFold(o, current) {
			return true
		}
	}
	return false
}
