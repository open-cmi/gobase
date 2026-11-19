package confparser

import (
	"encoding/json"
	"errors"
)

type Option struct {
	Name      string
	ParseFunc func(json.RawMessage) error
	SaveFunc  func() json.RawMessage
}

// Context config context
type Context struct {
	parser  *Parser
	Conf    map[string]json.RawMessage
	options map[string]Option
	order   []string
}

func NewContext() *Context {
	ctx := new(Context)
	ctx.Conf = make(map[string]json.RawMessage)
	ctx.options = make(map[string]Option)
	return ctx
}

func (c *Context) Register(opt *Option) error {
	_, found := c.options[opt.Name]
	if found {
		return errors.New("config " + opt.Name + " has been registered")
	}
	if opt.ParseFunc == nil {
		return errors.New("init func should not be empty")
	}

	if opt.SaveFunc == nil {
		return errors.New("save func should not be empty")
	}
	c.options[opt.Name] = *opt
	c.order = append(c.order, opt.Name)
	return nil
}

func (c *Context) Load(configfile string) error {
	if c.parser == nil {
		parser := NewParser(configfile)
		if parser == nil {
			return errors.New("open file failed")
		}
		c.parser = parser
	}

	err := c.parser.Load(&c.Conf)
	if err != nil {
		return err
	}

	for _, name := range c.order {
		option := c.options[name]
		value, ok := c.Conf[name]
		if ok {
			if option.ParseFunc != nil {
				err := option.ParseFunc(value)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Save save config
func (c *Context) Save() {
	for _, name := range c.order {
		option := c.options[name]
		if option.SaveFunc != nil {
			str := option.SaveFunc()
			c.Conf[name] = str
		}
	}
	if c.parser != nil {
		c.parser.Save(c.Conf)
	}
}
