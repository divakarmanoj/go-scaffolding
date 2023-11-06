# go-scaffolding

This is a scaffolding tool for go that generated go HTTP Rest API code given a model configuration. The configuration is accepted in YML file or through a interactive command line interface

Usage:
  go run main.go [command]

Available Commands:
  - completion  Generate the autocompletion script for the specified shell
  - generate    A brief description of your command
  - help        Help about any command
  - interactive This command will start an interactive command line interface to generate the model

Flags:
  -h, --help                help for go-scaffolding
      --output-dir string   The directory where the generated code will be placed
  -t, --toggle              Help message for toggle

Use "go run main.go [command] --help" for more information about a command.

## Step-by-Step Demonstration

1. Start by installing go-scaffolding. You can do this by running `go install` or by cloning the repository and building it manually.
2. Next, create a configuration file using the example provided in the link below:
```yaml
---
name: person
attributes:
  - name: name
    type: string
    is_required: false
  - name: age
    type: int16
    is_required: true
  - name: address
    type: struct
    attributes:
      - name: street_name
        type: string
        is_required: true
      - name: city
        type: string
        is_required: true
      - name: state
        type: string
        is_required: false
      - name: zip
        type: int16
        is_required: true

```
3. After creating the config file, run the following command in your terminal:
    
    ```bash
    go-scaffolding generate --config-file ./example/config.yaml
    ```
    
4. This command will generate the code in a new folder within the current directory. if you want to write code to a different directory use the command line arg `--output-dir`
