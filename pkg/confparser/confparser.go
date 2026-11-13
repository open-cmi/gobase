package confparser

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

// Parser parser struct
type Parser struct {
	path string
	fd   *os.File
}

// NewParser create a new parser
func NewParser(file string) *Parser {
	if file == "" || !strings.HasSuffix(file, ".json") {
		return nil
	}

	fd, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return nil
	}
	return &Parser{
		path: file,
		fd:   fd,
	}
}

// Load load data from file description
func (p *Parser) Load(v interface{}) error {
	p.fd.Seek(0, 0)

	content, err := io.ReadAll(p.fd)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}

// Save save data to file
func (p *Parser) Save(v interface{}) error {
	p.fd.Seek(0, 0)

	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	p.fd.Truncate(0)
	_, err = p.fd.Write(content)
	return err
}
