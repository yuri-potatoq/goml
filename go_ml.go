package go_ml

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

/* HTML element definitions */
type ElementType string

const (
	NonVoid ElementType = "non-void"
	Void    ElementType = "void"
)

type ContentType string

const (
	Raw  ContentType = "raw-text"
	Node ContentType = "node"
)

type HTMLRawContent struct {
	text string
}

type HTMLElement struct {
	tagName  string
	attrs    []HTMLAttribute
	contents []HTMLContent
	elType   ElementType
}

type HTMLContent struct {
	ctType ContentType
	child  HTMLElement
	raw    HTMLRawContent
}

// ref: https://github.com/golang/go/issues/62005#issuecomment-1747630201
type NopWritter struct{}

func (NopWritter) Write([]byte) (int, error) { return 0, nil }

func NopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(new(NopWritter), nil))
}

type buildConfig struct {
	stdWriter io.Writer
	debug     struct {
		logger *slog.Logger
	}
	indentitation struct {
		isEnable         bool
		indentationLevel uint8
	}
}

type buildOpt func(config *buildConfig)

func WithWriter(w io.Writer) buildOpt {
	return func(config *buildConfig) {
		config.stdWriter = w
	}
}

// TODO: create logger interface
func WithLogger(logger *slog.Logger) buildOpt {
	return func(config *buildConfig) {
		config.debug.logger = logger
	}
}

func WithDefaultIndentation() buildOpt {
	return func(config *buildConfig) {
		config.indentitation.indentationLevel = 4
		config.indentitation.isEnable = true
	}
}

// Choose a indentation level between 0 and 255.
func WithIndentation(level int) buildOpt {
	return func(config *buildConfig) {
		config.indentitation.indentationLevel = uint8(level)
		config.indentitation.isEnable = true
	}
}

func (ele HTMLElement) BuildDOM(opts ...buildOpt) error {
	return buildDOM(ele, opts...)
}

func (ct HTMLContent) BuildDOM(opts ...buildOpt) error {
	return buildDOM(ct.child, opts...)
}

func buildDOM(ele HTMLElement, opts ...buildOpt) error {
	var totalWritten int
	defaultCfg := new(buildConfig)
	defaultCfg.debug.logger = NopLogger()
	for _, op := range opts {
		op(defaultCfg)
	}

	if defaultCfg.stdWriter == nil {
		return errors.New("writer not found!")
	}

	// hardcoded html document compliance
	if ele.tagName == "html" {
		totalWritten, _ = defaultCfg.stdWriter.Write([]byte("<!DOCTYPE html>"))
	}

	n, err := defaultCfg.parseElement(ele, 1)
	if err != nil {
		return err
	}
	totalWritten += n

	defaultCfg.debug.logger.Info("Element was written with: ",
		"tag_name", ele.tagName, "bytes", totalWritten)
	return nil
}

func (cfg *buildConfig) parseElement(ele HTMLElement, tagDepth int) (int, error) {
	var attrStr string
	var rIndentStr, lIndentStr string
	var attrKeys []string
	var totalWritten int

	writeOrErr := func(s string) error {
		n, err := cfg.stdWriter.Write([]byte(s))
		totalWritten += n
		return err
	}

	putNChar := func(s string, ch string, n int) string {
		for i := 0; i < n; i++ {
			s += ch
		}
		return s
	}

	if cfg.indentitation.isEnable {
		rIndentStr = putNChar("\n", " ", tagDepth*int(cfg.indentitation.indentationLevel))
		lIndentStr = putNChar("\n", " ", (tagDepth-1)*int(cfg.indentitation.indentationLevel))
	}

	// rules:
	// 1. we need to merge all attributes with the same name
	// 2. if element has no attrs, the space become a suffix and will be removed
	attrMap := make(map[string]HTMLAttribute)
	for _, attr := range ele.attrs {
		if curAttr, ok := attrMap[attr.name]; ok {
			curAttr.values = append(curAttr.values, attr.values...)
			attrMap[attr.name] = curAttr
		} else {
			attrKeys = append(attrKeys, attr.name)
			attrMap[attr.name] = attr
		}
	}

	// TODO: find another aproach to have all the parsed attributes in O(n)
	for _, k := range attrKeys {
		attrStr += " " + attrMap[k].String()
	}
	attrStr = strings.TrimSuffix(attrStr, " ")

	switch ele.elType {
	// void -> <[tag][?attrs]/>
	case Void:
		if err := writeOrErr(fmt.Sprintf(`<%s%s/>`, ele.tagName, attrStr)); err != nil {
			return totalWritten, err
		}
		return totalWritten, nil

	// non-void -> <[tag][?attrs]>[content]</[tag]>
	default:
		if err := writeOrErr("<" + ele.tagName + attrStr + ">"); err != nil {
			return totalWritten, err
		}

		if len(ele.contents) > 0 {
			_ = writeOrErr(rIndentStr)
		}

		// threat as element node or just an raw text
		for _, ct := range ele.contents {
			switch ct.ctType {
			case Node:
				n, err := cfg.parseElement(ct.child, tagDepth+1)
				if err != nil {
					return totalWritten, err
				}
				totalWritten += n
			case Raw:
				if err := writeOrErr(ct.raw.text); err != nil {
					return totalWritten, err
				}
			default:
				return totalWritten, fmt.Errorf("not recognized content type: [%s]", ct.ctType)
			}
		}

		if len(ele.contents) > 0 {
			_ = writeOrErr(lIndentStr)
		}

		if err := writeOrErr("</" + ele.tagName + ">"); err != nil {
			return totalWritten, err
		}
		return totalWritten, nil
	}
}

/*
	HTML attribute definitions

Ref: https://www.w3.org/TR/2012/WD-html-markup-20120329/syntax.html#syntax-attributes
*/
type AttributeType string

const (
	None         AttributeType = "none"
	Single       AttributeType = "single"
	DoubleQuoted AttributeType = "double-quoted"
)

type HTMLAttribute struct {
	name string
	// space delimited values
	values   []string
	attrType AttributeType
}

func (attr HTMLAttribute) String() string {
	switch attr.attrType {
	case DoubleQuoted:
		var st string
		for _, v := range attr.values {
			st += v + " "
		}
		return fmt.Sprintf(`%s="%s"`, attr.name, strings.TrimSuffix(st, " "))
	case Single:
		return attr.name
	case None:
		return ""
	default:
		return ""
	}
}

/* Attributes functions declarations */
func Attr(name string, attrType AttributeType, values ...string) HTMLAttribute {
	return HTMLAttribute{name: name, values: values, attrType: attrType}
}

func ClassNames(values ...string) HTMLAttribute {
	return Attr("class", DoubleQuoted, values...)
}

func PlaceHolder(value string) HTMLAttribute {
	return Attr("placeholder", DoubleQuoted, value)
}

func Id(values ...string) HTMLAttribute {
	return Attr("id", DoubleQuoted, values...)
}

func Name(values ...string) HTMLAttribute {
	return Attr("name", DoubleQuoted, values...)
}

func Lang(values ...string) HTMLAttribute {
	return Attr("lang", DoubleQuoted, values...)
}

func Type(values ...string) HTMLAttribute {
	return Attr("type", DoubleQuoted, values...)
}

func Value(values ...string) HTMLAttribute {
	return Attr("value", DoubleQuoted, values...)
}

func Src(values ...string) HTMLAttribute {
	return Attr("src", DoubleQuoted, values...)
}

func Defer() HTMLAttribute {
	return Attr("defer", Single)
}

func Checked() HTMLAttribute {
	return Attr("checked", Single)
}

func Required() HTMLAttribute {
	return Attr("required", Single)
}

func Action(values ...string) HTMLAttribute {
	return Attr("action", Single, values...)
}

func Method(values ...string) HTMLAttribute {
	return Attr("method", Single, values...)
}

/* Attributes utils */
func IsChecked(check bool) (attr HTMLAttribute) {
	if check {
		return Checked()
	}
	return HTMLAttribute{attrType: None}
}

/* Tags functions declarations */
type tagClosure func(contents ...HTMLContent) HTMLContent

func Tag(tagName string, elType ElementType, attrs ...HTMLAttribute) tagClosure {
	return func(contents ...HTMLContent) HTMLContent {
		return HTMLContent{
			child: HTMLElement{
				tagName:  tagName,
				contents: contents,
				attrs:    attrs,
				elType:   elType,
			},
			ctType: Node,
		}
	}
}

func RawText(text string) HTMLContent {
	return HTMLContent{raw: HTMLRawContent{text: text}, ctType: Raw}
}

func Input(attrs ...HTMLAttribute) HTMLContent {
	return Tag("input", Void, attrs...)()
}

func Button(attrs ...HTMLAttribute) tagClosure {
	return Tag("button", NonVoid, attrs...)
}

func Script(attrs ...HTMLAttribute) tagClosure {
	return Tag("script", NonVoid, attrs...)
}

func Title(attrs ...HTMLAttribute) tagClosure {
	return Tag("title", NonVoid, attrs...)
}

func Head(attrs ...HTMLAttribute) tagClosure {
	return Tag("head", NonVoid, attrs...)
}

func Div(attrs ...HTMLAttribute) tagClosure {
	return Tag("div", NonVoid, attrs...)
}

func Body(attrs ...HTMLAttribute) tagClosure {
	return Tag("body", NonVoid, attrs...)
}

func Html(attrs ...HTMLAttribute) tagClosure {
	return Tag("html", NonVoid, attrs...)
}

func Form(attrs ...HTMLAttribute) tagClosure {
	return Tag("form", NonVoid, attrs...)
}

func Label(attrs ...HTMLAttribute) tagClosure {
	return Tag("label", NonVoid, attrs...)
}

func Table(attrs ...HTMLAttribute) tagClosure {
	return Tag("table", NonVoid, attrs...)
}

func Th(attrs ...HTMLAttribute) tagClosure {
	return Tag("th", NonVoid, attrs...)
}

func Tr(attrs ...HTMLAttribute) tagClosure {
	return Tag("tr", NonVoid, attrs...)
}

func Td(attrs ...HTMLAttribute) tagClosure {
	return Tag("td", NonVoid, attrs...)
}
