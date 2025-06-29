package i18n

import (
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printerMapping map[string]*message.Printer = make(map[string]*message.Printer)

var supportedLang []string = []string{"en-US", "zh-CN"}

var gLang string

// Sprint is like fmt.Sprint, but using language-specific formatting.
func Sprint(args ...interface{}) string {
	return printerMapping[gLang].Sprint(args)
}

// Print is like fmt.Print, but using language-specific formatting.
func Print(args ...interface{}) (n int, err error) {
	return printerMapping[gLang].Print(args)
}

// Sprintln is like fmt.Sprintln, but using language-specific formatting.
func Sprintln(args ...interface{}) string {
	return printerMapping[gLang].Sprintln(args)
}

// Println is like fmt.Println, but using language-specific formatting.
func Println(args ...interface{}) (n int, err error) {
	return printerMapping[gLang].Println(args)
}

// Sprintf is like fmt.Sprintf, but using language-specific formatting.
func Sprintf(format string, args ...interface{}) string {
	return printerMapping[gLang].Sprintf(format, args)
}

// Printf is like fmt.Printf, but using language-specific formatting.
func Printf(format string, args ...interface{}) (n int, err error) {
	return printerMapping[gLang].Printf(format, args)
}

func SetLang(lang string) error {
	var support bool = false
	for _, l := range supportedLang {
		if l == lang {
			support = true
		}
	}
	if !support {
		return errors.New("lang not support")
	}
	gLang = lang
	return nil
}

func GetLang() string {
	return gLang
}

func init() {
	gLang = "zh-CN"
	for _, lang := range supportedLang {
		tag := language.MustParse(lang)
		p := message.NewPrinter(tag)
		printerMapping[lang] = p
	}
}
