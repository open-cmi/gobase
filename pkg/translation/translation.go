package translation

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PlaceHolder struct {
	ID             string `json:"id"`
	String         string `json:"string"`
	ArgNum         int    `json:"argNum"`
	UnderlyingType string `json:"underlyingType"`
}

type MessageItem struct {
	ID           string        `json:"id"`
	Message      string        `json:"message"`
	Translation  string        `json:"translation"`
	PlaceHolders []PlaceHolder `json:"placeholders"`
}

type MessageStruct struct {
	Language string        `json:"language"`
	Messages []MessageItem `json:"messages"`
}

type Dictionary struct {
	trans map[string]string
}

func (d *Dictionary) Lookup(key string) (data string, ok bool) {
	trans, ok := d.trans[key]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("\x02%s", trans), true
}

func InitTransDict(content string) (*Dictionary, error) {
	d := &Dictionary{
		trans: make(map[string]string),
	}

	var m MessageStruct
	err := json.Unmarshal([]byte(content), &m)
	if err != nil {
		return nil, err
	}

	for _, item := range m.Messages {
		if len(item.PlaceHolders) != 0 {
			for _, ph := range item.PlaceHolders {
				item.Translation = strings.Replace(item.Translation, fmt.Sprintf("{%s}", ph.ID), ph.String, 1)

				placeHolder := fmt.Sprintf("[%d]", ph.ArgNum)
				item.ID = strings.Replace(item.ID, fmt.Sprintf("{%s}", ph.ID), strings.ReplaceAll(ph.String, placeHolder, ""), -1)
			}
		}
		d.trans[item.ID] = item.Translation
	}
	return d, nil
}
