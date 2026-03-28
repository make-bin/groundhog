package skill

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// Frontmatter holds parsed SKILL.md frontmatter fields.
type Frontmatter struct {
	Name                   string
	Description            string
	Always                 bool     // always inject into system prompt
	UserInvocable          bool     // show in user-facing skill list (default true)
	DisableModelInvocation bool     // exclude from system prompt injection
	Requires               []string // binary names that must exist on PATH
	OS                     []string // platform filter: "darwin", "linux", "windows"
}

// parseFrontmatterFromFile reads a SKILL.md file and returns its frontmatter and body.
func parseFrontmatterFromFile(path string) (Frontmatter, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return Frontmatter{}, "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var fmLines []string
	var bodyLines []string
	inFM := false
	fmDone := false

	for scanner.Scan() {
		line := scanner.Text()
		if !fmDone {
			if line == "---" {
				if !inFM {
					inFM = true
					continue
				}
				fmDone = true
				continue
			}
			if inFM {
				fmLines = append(fmLines, line)
				continue
			}
			// no frontmatter opener found
			fmDone = true
		}
		bodyLines = append(bodyLines, line)
	}
	if err := scanner.Err(); err != nil {
		return Frontmatter{}, "", err
	}

	fm := parseFrontmatterLines(fmLines)
	body := strings.Join(bodyLines, "\n")
	return fm, body, nil
}

// parseFrontmatterLines parses simple YAML key: value lines.
func parseFrontmatterLines(lines []string) Frontmatter {
	fm := Frontmatter{
		UserInvocable: true, // default true per openclaw spec
	}
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "name":
			fm.Name = val
		case "description":
			fm.Description = val
		case "always":
			fm.Always = parseBool(val, false)
		case "user-invocable":
			fm.UserInvocable = parseBool(val, true)
		case "disable-model-invocation":
			fm.DisableModelInvocation = parseBool(val, false)
		case "requires":
			fm.Requires = parseStringList(val)
		case "os":
			fm.OS = parseStringList(val)
		}
	}
	return fm
}

func parseBool(s string, defaultVal bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "1":
		return true
	case "false", "no", "0":
		return false
	}
	return defaultVal
}

// parseStringList parses "a, b, c" or "[a, b, c]" into a slice.
func parseStringList(s string) []string {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// hasBinary checks whether a binary exists on PATH.
func hasBinary(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
