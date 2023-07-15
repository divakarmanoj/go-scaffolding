package generator

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"unicode"
)

// ParseStruct function accepts a string and tries to parse it into a Struct
func ParseStruct(s string) (*Config, error) {
	var Data Config

	err := yaml.Unmarshal([]byte(s), &Data)
	if err != nil {
		err := json.Unmarshal([]byte(s), &Data)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return &Data, nil
}

func ParseStructFromFileName(fileName string) (*Config, error) {
	// Read the file
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseStruct(string(dat))
}

func toCamelCase(name string) string {
	words := strings.Split(name, "_")
	result := ""
	if name == "" {
		return name
	}
	for _, word := range words {
		result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}

	return result
}

func toTitleCase(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		s)
}

func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

func ToSnakeCase(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.ToLower(name)

	return name
}

func isValidType(fieldType string) bool {
	validTypes := []string{
		"string",
		"int64",
		"int32",
		"int16",
		"float64",
		"bool",
		"time.Time",
		"struct",
	}
	for _, t := range validTypes {
		if fieldType == t {
			return true
		}
	}
	return false
}
