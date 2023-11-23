package go_ml

import (
	"fmt"
	"strings"
)

/* HTML element definitions */
type ElementType string

const (
	NonVoid ElementType = "non-void"
	Void    ElementType = "void"
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
	child HTMLElement
	raw   HTMLRawContent
}

func (ele HTMLElement) String() string {
	var contentStr strings.Builder
	var attrStr string

	// if element has no attrs, the space become a suffix and will be removed
	attrStr += " "
	for _, attr := range ele.attrs {
		attrStr += attr.String()
	}

	switch ele.elType {
	// void -> <[tag][?attrs]/>
	case Void:
		return fmt.Sprintf(`<%s%s/>`, ele.tagName, strings.TrimSuffix(attrStr, " "))

	// non-void -> <[tag][?attrs]>[content]</[tag]>
	default:
		// treat as element node or just an raw text
		for _, ct := range ele.contents {
			// TODO: handle write error correctly
			if ct.raw.text != "" {
				contentStr.WriteString(ct.raw.text)
			} else {
				contentStr.WriteString(ct.child.String())
			}

		}
		return fmt.Sprintf(`<%s%s>%s</%s>`,
			ele.tagName, strings.TrimSuffix(attrStr, " "), contentStr.String(), ele.tagName)
	}
}

func (ele HTMLElement) BuildDOM() string {
	return buildDOM(ele)
}

func (ct HTMLContent) BuildDOM() string {
	return buildDOM(ct.child)
}

func buildDOM(ele HTMLElement) string {
	var main strings.Builder

	// hardcoded html document compliance
	if ele.tagName == "html" {
		main.WriteString("<!DOCTYPE html>")
	}

	// TODO: handle write error correctly
	main.WriteString(ele.String())
	return main.String()
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
		}
	}
}

func RawText(text string) HTMLContent {
	return HTMLContent{raw: HTMLRawContent{text: text}}
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
