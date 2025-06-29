package confparser

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser parser struct
type Parser struct {
	path   string
	fd     *os.File
	format string
}

// NewParser create a new parser
func NewParser(file string) *Parser {
	if file == "" {
		return nil
	}

	var format string = ""
	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		format = "yaml"
	} else if strings.HasSuffix(file, "json") {
		format = "json"
	} else {
		return nil
	}
	fd, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return nil
	}
	return &Parser{
		path:   file,
		fd:     fd,
		format: format,
	}
}

// Load load data from file description
func (p *Parser) Load(v interface{}) (err error) {
	p.fd.Seek(0, 0)

	content, err := io.ReadAll(p.fd)
	if err != nil {
		return err
	}
	switch p.format {
	case "json":
		return json.Unmarshal(content, v)
	case "yaml":
		return yaml.Unmarshal(content, v)
	default:
		return errors.New("file format not suppport")
	}
}

// Save save data to file
func (p *Parser) Save(v interface{}) (err error) {
	p.fd.Seek(0, 0)

	var content []byte
	switch p.format {
	case "json":
		content, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
	case "yaml":
		content, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	default:
		return errors.New("file format not suppport")
	}

	p.fd.Truncate(0)
	_, err = p.fd.Write(content)

	return err
}
