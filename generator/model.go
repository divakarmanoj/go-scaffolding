package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"unicode"
)

type Attribute struct {
	Name       string `json:"name" yaml:"name"`
	camelCase  string
	Type       string      `json:"type" yaml:"type"`
	Attributes []Attribute `json:"attributes" yaml:"attributes"`
	IsRequired bool        `json:"is_required" yaml:"is_required"`
}

type Config struct {
	Name       string `json:"name" yaml:"name"`
	camelCase  string
	Attributes []Attribute `json:"attributes" yaml:"attributes"`
}

func (receiver *Config) ToYAML() (string, error) {
	marshal, err := yaml.Marshal(receiver)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (receiver *Config) ToJSON() (string, error) {
	marshal, err := json.Marshal(receiver)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (receiver *Config) Generate(outputDir string) {
	err := receiver.Validate()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	if outputDir == "" {
		outputDir = "."
	}
	if rune(outputDir[len(outputDir)-1]) != os.PathSeparator {
		outputDir += "/"
	}
	receiver.camelCase = toCamelCase(receiver.Name)
	for i, attr := range receiver.Attributes {
		receiver.Attributes[i] = AttributeCamelCase(attr)
	}
	fmt.Printf("%+v\n", receiver)
	// Create directory if not exists with name of struct
	_ = os.MkdirAll(outputDir+ToSnakeCase(receiver.Name), os.ModePerm)
	GenerateRequestResponse(receiver, outputDir)
	_, model := GenerateModel(receiver, outputDir)
	GenerateHandler(receiver, outputDir)
	GenerateMain(receiver, model, outputDir)
}

func AttributeCamelCase(attr Attribute) Attribute {
	attr.camelCase = toCamelCase(attr.Name)
	if attr.Type == "struct" {
		for i, a := range attr.Attributes {
			attr.Attributes[i] = AttributeCamelCase(a)
		}
	}
	return attr
}

func (receiver *Config) Validate() error {
	if receiver.Name == "" {
		return errors.New("name is required")
	}
	if err := NameValaidation(receiver.Name); err != nil {
		return err
	}
	if len(receiver.Attributes) == 0 {
		return errors.New("attributes is required")
	}
	for _, attr := range receiver.Attributes {
		err := attr.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (receiver *Attribute) Validate() error {
	if receiver.Name == "" {
		return errors.New("name is required")
	}
	if err := NameValaidation(receiver.Name); err != nil {
		return err
	}
	if !isValidType(receiver.Type) {
		return fmt.Errorf("invalid type %s", receiver.Type)
	}
	if receiver.Type == "struct" && len(receiver.Attributes) == 0 {
		return errors.New("attributes is required")
	}
	if receiver.Type == "struct" {
		for _, attr := range receiver.Attributes {
			err := attr.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NameValaidation(input string) error {
	if len(input) < 1 {
		return errors.New("name must be at least 1 character")
	}
	// first character must be a letter
	if unicode.IsNumber(rune(input[0])) {
		return errors.New("first character must be a letter")
	}

	// only letters, numbers and underscores allowed
	for _, c := range input {
		if !unicode.IsLower(c) && !unicode.IsNumber(c) && c != '_' {
			return errors.New("only letters, numbers and underscores allowed")
		}
	}
	return nil
}
