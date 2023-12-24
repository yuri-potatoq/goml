package go_ml

import (
	"errors"
	xnet_html "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log/slog"
)

var (
	ErrWriterNotFound = errors.New("writer not found!")
)

/* Type aliases */
type (
	HtmlAttribute xnet_html.Attribute
	HtmlNode      struct{ *xnet_html.Node }
	HtmlNodeType  xnet_html.NodeType
)

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

func BuildDOM(w io.Writer, n *HtmlNode) error {
	return xnet_html.Render(w, n.Node)
}

/* Attributes functions declarations */
func Attr(at atom.Atom, values string) HtmlAttribute {
	return HtmlAttribute{Key: at.String(), Val: values}
}

func ClassNames(values string) HtmlAttribute {
	return Attr(atom.Class, values)
}

func PlaceHolder(value string) HtmlAttribute {
	return Attr(atom.Placeholder, value)
}

func Id(value string) HtmlAttribute {
	return Attr(atom.Id, value)
}

func Name(value string) HtmlAttribute {
	return Attr(atom.Name, value)
}

func Lang(value string) HtmlAttribute {
	return Attr(atom.Lang, value)
}

func Type(value string) HtmlAttribute {
	return Attr(atom.Type, value)
}

func Value(value string) HtmlAttribute {
	return Attr(atom.Value, value)
}

func Src(value string) HtmlAttribute {
	return Attr(atom.Src, value)
}

func Defer() HtmlAttribute {
	return Attr(atom.Defer, "")
}

func Checked() HtmlAttribute {
	return Attr(atom.Checked, "")
}

func Required() HtmlAttribute {
	return Attr(atom.Required, "")
}

func Action(value string) HtmlAttribute {
	return Attr(atom.Action, value)
}

func Method(value string) HtmlAttribute {
	return Attr(atom.Method, value)
}

/* Attributes utils */
func IsChecked(check bool) HtmlAttribute {
	if check {
		return Checked()
	}
	return HtmlAttribute{}
}

type tagClosure func(contents ...*HtmlNode) *HtmlNode

func Tag(tagname atom.Atom, nt xnet_html.NodeType, attrs ...HtmlAttribute) tagClosure {
	return func(contents ...*HtmlNode) *HtmlNode {
		n := &xnet_html.Node{
			Type: xnet_html.NodeType(nt),
			Data: tagname.String(),
		}

		for _, attr := range attrs {
			n.Attr = append(n.Attr, xnet_html.Attribute(attr))
		}

		for _, ct := range contents {
			// TODO: add indentation feature
			// n.InsertBefore(&net_html.Node{Data: "\n\t", Type: net_html.TextNode}, nil)
			n.InsertBefore(ct.Node, nil)
			// n.InsertBefore(&net_html.Node{Data: "\n", Type: net_html.TextNode}, nil)
		}
		return &HtmlNode{Node: n}
	}
}

func Input(attrs ...HtmlAttribute) *HtmlNode {
	return Tag(atom.Input, xnet_html.ElementNode, attrs...)()
}

func Div(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Div, xnet_html.ElementNode, attrs...)
}

func Comment(text string) *HtmlNode {
	return &HtmlNode{&xnet_html.Node{Type: xnet_html.CommentNode, Data: text}}
}

func RawText(text string) *HtmlNode {
	return &HtmlNode{&xnet_html.Node{Type: xnet_html.RawNode, Data: text}}
}

func Button(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Button, xnet_html.ElementNode, attrs...)
}

func Script(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Script, xnet_html.ElementNode, attrs...)
}

func Title(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Title, xnet_html.ElementNode, attrs...)
}

func Head(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Head, xnet_html.ElementNode, attrs...)
}

func Body(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Body, xnet_html.ElementNode, attrs...)
}

func Html(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Html, xnet_html.ElementNode, attrs...)
}

func Form(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Form, xnet_html.ElementNode, attrs...)
}

func Label(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Label, xnet_html.ElementNode, attrs...)
}

func Table(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Table, xnet_html.ElementNode, attrs...)
}

func Th(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Th, xnet_html.ElementNode, attrs...)
}

func Tr(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Tr, xnet_html.ElementNode, attrs...)
}

func Td(attrs ...HtmlAttribute) tagClosure {
	return Tag(atom.Td, xnet_html.ElementNode, attrs...)
}
