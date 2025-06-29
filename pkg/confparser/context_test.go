package confparser

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConfig1(t *testing.T) {

	ctx := NewContext()
	if ctx == nil {
		t.Errorf("NewContext failed\n")
	}

	conf := make(map[string]interface{}, 0)

	master := make(map[string]interface{})
	master["master"] = "http://localhost:30016"
	master["useunix"] = false
	conf["transport2"] = master

	contentByte, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		t.Errorf("json MarshalIndent failed\n")
	}
	testFile := "/tmp/config_test.json"
	fh, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("open config_test file failed\n")
	}
	fh.WriteString(string(contentByte))

	var opt Option
	opt.Name = "transport2"
	opt.ParseFunc = func(raw json.RawMessage) error {
		return nil
	}

	opt.SaveFunc = func() json.RawMessage {
		v, _ := json.Marshal(master)
		return v
	}

	err = ctx.Register(&opt)
	if err != nil {
		t.Errorf("ctx Register failed: %s\n", err.Error())
	}

	err = ctx.Load(testFile)
	if err != nil {
		t.Errorf("ctx Load file failed: %s\n", err.Error())
	}
	ctx.Save()
	os.Remove(testFile)
}
