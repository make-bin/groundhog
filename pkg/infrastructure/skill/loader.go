// Package skill implements openclaw-compatible skill loading.
//
// Skills are Markdown files (SKILL.md) with optional YAML frontmatter.
// This package mirrors the openclaw TypeScript skill loading logic:
//   - Scans multiple directories in priority order
//   - Later sources override earlier ones by skill name
//   - Filters by OS, binary requirements, and session skill list
//   - Builds a system prompt from eligible skills
package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	maxSkillFileBytes  = 100 * 1024 // 100 KB per SKILL.md
	maxSkillsPerSource = 50
	skillFileName      = "SKILL.md"
)

// Registry holds all loaded skills merged from multiple sources.
type Registry struct {
	entries []*Entry // ordered by name for determinism
	byName  map[string]*Entry
}

// SourceConfig defines a skill source directory and its label.
type SourceConfig struct {
	Dir    string
	Source string
}

// LoadConfig controls which directories are scanned.
type LoadConfig struct {
	// BundledDir is the built-in skills shipped with the binary (lowest priority).
	BundledDir string
	// ManagedDir is ~/.openclaw/skills (user-installed skills).
	ManagedDir string
	// PersonalAgentsDir is ~/.agents/skills.
	PersonalAgentsDir string
	// ProjectAgentsDir is <workspaceDir>/.agents/skills.
	ProjectAgentsDir string
	// WorkspaceDir is <workspaceDir>/skills (highest priority).
	WorkspaceDir string
	// ExtraDirs are additional directories from config (between bundled and managed).
	ExtraDirs []string
}

// NewRegistry loads skills from all configured directories.
// Priority (lowest → highest, later overrides earlier by name):
//
//	extra dirs < bundled < managed < ~/.agents/skills < .agents/skills < workspace/skills
func NewRegistry(cfg LoadConfig) *Registry {
	r := &Registry{byName: make(map[string]*Entry)}

	sources := buildSourceList(cfg)
	for _, src := range sources {
		entries := scanDir(src.Dir, src.Source)
		for _, e := range entries {
			r.byName[e.Name] = e
		}
	}

	// Build sorted slice for deterministic iteration
	r.entries = make([]*Entry, 0, len(r.byName))
	for _, e := range r.byName {
		r.entries = append(r.entries, e)
	}
	sort.Slice(r.entries, func(i, j int) bool {
		return r.entries[i].Name < r.entries[j].Name
	})

	return r
}

// List returns all loaded entries sorted by name.
func (r *Registry) List() []*Entry {
	result := make([]*Entry, len(r.entries))
	copy(result, r.entries)
	return result
}

// Get returns the entry for the given skill name, or nil.
func (r *Registry) Get(name string) *Entry {
	return r.byName[name]
}

// ResolvePrompt returns the system prompt fragment for all eligible skills.
// workspaceSkillNames: if non-empty, only those skills (plus always=true) are included.
// This is called fresh on every agent run so new/removed SKILL.md files take effect immediately.
func (r *Registry) ResolvePrompt(workspaceSkillNames []string) string {
	ctx := EligibilityContext{WorkspaceSkillNames: workspaceSkillNames}

	var eligible []*Entry
	for _, e := range r.entries {
		if e.IsEligible(ctx) && e.ShouldInjectIntoPrompt() {
			eligible = append(eligible, e)
		}
	}

	if len(eligible) == 0 {
		return ""
	}

	return formatSkillsPrompt(eligible)
}

// formatSkillsPrompt renders eligible skills as XML blocks, matching openclaw TS format.
func formatSkillsPrompt(entries []*Entry) string {
	var sb strings.Builder
	sb.WriteString("<skills>\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("<skill name=%q", e.Name))
		if e.Description != "" {
			sb.WriteString(fmt.Sprintf(" description=%q", e.Description))
		}
		sb.WriteString(">\n")
		sb.WriteString(strings.TrimSpace(e.Content))
		sb.WriteString("\n</skill>\n")
	}
	sb.WriteString("</skills>")
	return sb.String()
}

// buildSourceList returns sources in priority order (lowest first).
func buildSourceList(cfg LoadConfig) []SourceConfig {
	var sources []SourceConfig

	for _, d := range cfg.ExtraDirs {
		if d != "" {
			sources = append(sources, SourceConfig{Dir: d, Source: "extra"})
		}
	}
	if cfg.BundledDir != "" {
		sources = append(sources, SourceConfig{Dir: cfg.BundledDir, Source: "bundled"})
	}
	if cfg.ManagedDir != "" {
		sources = append(sources, SourceConfig{Dir: cfg.ManagedDir, Source: "managed"})
	}
	if cfg.PersonalAgentsDir != "" {
		sources = append(sources, SourceConfig{Dir: cfg.PersonalAgentsDir, Source: "agents-personal"})
	}
	if cfg.ProjectAgentsDir != "" {
		sources = append(sources, SourceConfig{Dir: cfg.ProjectAgentsDir, Source: "agents-project"})
	}
	if cfg.WorkspaceDir != "" {
		sources = append(sources, SourceConfig{Dir: cfg.WorkspaceDir, Source: "workspace"})
	}
	return sources
}

// scanDir scans a single directory for skill subdirectories containing SKILL.md.
func scanDir(dir, source string) []*Entry {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // directory doesn't exist or unreadable — skip silently
	}

	var skills []*Entry
	for _, de := range entries {
		if !de.IsDir() {
			continue
		}
		skillPath := filepath.Join(dir, de.Name(), skillFileName)
		e := loadSkillFile(skillPath, de.Name(), source)
		if e == nil {
			continue
		}
		skills = append(skills, e)
		if len(skills) >= maxSkillsPerSource {
			break
		}
	}
	return skills
}

// loadSkillFile reads and parses a single SKILL.md file.
// Returns nil if the file doesn't exist, is too large, or fails to parse.
func loadSkillFile(path, dirName, source string) *Entry {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if info.Size() > maxSkillFileBytes {
		return nil // skip oversized files
	}

	fm, body, err := parseFrontmatterFromFile(path)
	if err != nil {
		return nil
	}

	name := fm.Name
	if name == "" {
		name = dirName
	}

	return &Entry{
		Name:        name,
		Description: fm.Description,
		Content:     body,
		Frontmatter: fm,
		Source:      source,
		FilePath:    path,
	}
}

// DefaultLoadConfig builds a LoadConfig from the given workspace directory.
// Mirrors openclaw TS directory resolution.
func DefaultLoadConfig(workspaceDir string, extraDirs []string) LoadConfig {
	home, _ := os.UserHomeDir()
	return LoadConfig{
		ManagedDir:        filepath.Join(home, ".openclaw", "skills"),
		PersonalAgentsDir: filepath.Join(home, ".agents", "skills"),
		ProjectAgentsDir:  filepath.Join(workspaceDir, ".agents", "skills"),
		WorkspaceDir:      filepath.Join(workspaceDir, "skills"),
		ExtraDirs:         extraDirs,
	}
}
