package main

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

func (el HTMLElement) WithText(text string) HTMLElement {
	return HTMLElement{
		tagName:  el.tagName,
		attrs:    el.attrs,
		contents: el.contents,
		elType:   el.elType,
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

/* HTML attribute definitions
Ref: https://www.w3.org/TR/2012/WD-html-markup-20120329/syntax.html#syntax-attributes
*/

type HTMLAttribute struct {
	name string
	// space delimited values
	values []string
}

func (attr HTMLAttribute) String() string {
	var st string

	if attr.name != "" && len(attr.values) < 1 {
		return attr.name
	}

	for _, v := range attr.values {
		st += v + " "
	}
	return fmt.Sprintf(`%s="%s"`, attr.name, strings.TrimSuffix(st, " "))
}

/* Attributes functions declarations */
func attr(name string, values ...string) HTMLAttribute {
	return HTMLAttribute{name: name, values: values}
}

func ClassNames(values ...string) HTMLAttribute {
	return attr("class", values...)
}

func Lang(values ...string) HTMLAttribute {
	return attr("lang", values...)
}

func Type(values ...string) HTMLAttribute {
	return attr("type", values...)
}

func Defer() HTMLAttribute {
	return attr("defer")
}

func Required() HTMLAttribute {
	return attr("required")
}

/* Tags functions declarations */
type tagClosure func(contents ...HTMLContent) HTMLContent

func tag(tagName string, elType ElementType, attrs ...HTMLAttribute) tagClosure {
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
	return tag("input", Void, attrs...)()
}

func Script(attrs ...HTMLAttribute) tagClosure {
	return tag("script", NonVoid, attrs...)
}

func Head(attrs ...HTMLAttribute) tagClosure {
	return tag("head", NonVoid, attrs...)
}

func Div(attrs ...HTMLAttribute) tagClosure {
	return tag("div", NonVoid, attrs...)
}

func Body(attrs ...HTMLAttribute) tagClosure {
	return tag("body", NonVoid, attrs...)
}

func Html(attrs ...HTMLAttribute) tagClosure {
	return tag("html", NonVoid, attrs...)
}
