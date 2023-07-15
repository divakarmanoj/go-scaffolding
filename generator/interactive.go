package generator

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
)

func InteractiveConfigGeneration(modelName string, skipName bool) (*Config, error) {
	var err error
	config := &Config{}

	if !skipName {
		p := promptui.Prompt{
			Label:    "Model Name",
			Validate: NameValaidation,
		}

		modelName, err = p.Run()
		if err != nil {
			return nil, err
		}

		config.Name = modelName
	}
	fmt.Printf("Enter the attributes for the %s model\n", modelName)
	for true {
		attribute := &Attribute{}
		prompt := promptui.Prompt{
			Label:    "Attribute Name",
			Validate: NameValaidation,
		}
		attributeName, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		attribute.Name = attributeName
		datatypePrompt := promptui.Select{
			Label: "Select Datatype",
			Items: []string{"string", "int16", "int32", "int64", "float64", "bool", "time.Time", "struct"},
		}

		_, datatype, err := datatypePrompt.Run()
		if err != nil {
			return nil, err
		}

		attribute.Type = datatype

		isRequiredpPrompt := promptui.Select{
			Label: "Is this attribute required?",
			Items: []string{"true", "false"},
		}

		_, isRequired, err := isRequiredpPrompt.Run()
		if err != nil {
			return nil, err
		}

		attribute.IsRequired = isRequired == "true"

		config.Attributes = append(config.Attributes, *attribute)

		isAnotherAttributePrompt := promptui.Select{
			Label: "Do you want to add another attribute?",
			Items: []string{"true", "false"},
		}

		_, isAnotherAttribute, err := isAnotherAttributePrompt.Run()
		if err != nil {
			return nil, err
		}

		if isAnotherAttribute != "true" {
			break
		}
	}

	// if one of the attributes is a struct, we need to add the struct definition
	for i, attribute := range config.Attributes {
		if attribute.Type == "struct" {
			println()
			model, err := InteractiveConfigGeneration(attribute.Name, true)
			if err != nil {
				return nil, err
			}
			config.Attributes[i].Attributes = model.Attributes
		}
	}

	// print the config in yaml format and ask for confirmation
	if !skipName {
		fmt.Println("The model config is as follows:")
		fmt.Println(config.ToYAML())
		confirmPrompt := promptui.Select{
			Label: "Do you want to confirm the model?",
			Items: []string{"true", "false"},
		}
		if _, confirm, err := confirmPrompt.Run(); err != nil || confirm != "true" {
			if confirm != "true" {
				println("Aborted")
				return nil, errors.New("aborted")
			}
			return nil, err
		}
	}

	return config, nil
}
