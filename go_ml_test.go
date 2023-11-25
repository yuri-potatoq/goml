package go_ml

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestBasicDOM(t *testing.T) {
	testSuite := []struct {
		name         string
		givenDOM     HTMLContent
		expectedHtml string
	}{
		{
			name:         "build empty html document",
			expectedHtml: `<!DOCTYPE html><html lang="en"></html>`,
			givenDOM:     Html(Lang("en"))(),
		},
		{
			name:         "build bare div with class names",
			expectedHtml: `<div class="main container"><div></div></div>`,
			givenDOM:     Div(ClassNames("main", "container"))(Div()()),
		},
		{
			name:         "build script tag with non key-valued attributes",
			expectedHtml: `<script defer></script>`,
			givenDOM:     Script(Defer())(),
		},
		{
			name:         "build script tag with multiple attributes",
			expectedHtml: `<script defer src="./index.js"></script>`,
			givenDOM:     Script(Defer(), Src("./index.js"))(),
		},
		{
			name:         "build input with checked attribute",
			expectedHtml: `<input checked/>`,
			givenDOM:     Input(Checked()),
		},
		{
			name:         "build script tag with source files",
			expectedHtml: `<script src="index.js"></script>`,
			givenDOM:     Script(Src("index.js"))(),
		},
		{
			name:         "build input regular tag",
			expectedHtml: `<input type="text"/>`,
			givenDOM:     Input(Type("text")),
		},
		{
			name:         "build non-void tag with inner text",
			expectedHtml: `<div>olá &#10; mundo</div>`,
			givenDOM:     Div()(RawText("olá &#10; mundo")),
		},
		{
			name:         "build non-void tag with inner text",
			expectedHtml: `<div>olá &#10; mundo<div></div>test<div></div></div>`,
			givenDOM:     Div()(RawText("olá &#10; mundo"), Div()(), RawText("test"), Div()()),
		},
		{
			name:         "build input with mutiple class attributes",
			expectedHtml: `<input name="task" class="text name editable"/>`,
			givenDOM:     Input(Name("task"), ClassNames("text"), ClassNames("name"), ClassNames("editable")),
		},
		/* HTMX attributes tests */
		{
			name:         "build input with htmx attributes",
			expectedHtml: `<input checked hx-on:click="console.log('hello')"/>`,
			givenDOM:     Input(Checked(), HxOn("click", "console.log('hello')")),
		},
		// TODO: fix wrong attribute spaces sort
		// i.g.: <input type="checkbox"required required="required"/>
		// {
		// 	name:         "build all boolean attributes possibilities",
		// 	expectedHtml: `<input type="text" required required="required"/>`,
		// 	givenDOM:     Input(Type("checkbox"), Required(), attr("required", "required")),
		// },
	}

	debugLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	for _, tc := range testSuite {
		t.Run(tc.name, func(t *testing.T) {
			st := new(strings.Builder)
			err := tc.givenDOM.BuildDOM(WithLogger(debugLogger), WithWriter(st))
			if err != nil {
				t.Error(err)
			}
			parsedDoc := st.String()
			if strings.Compare(parsedDoc, tc.expectedHtml) != 0 {
				t.Errorf("result not match: given: [%s], expected: [%s]", parsedDoc, tc.expectedHtml)
			}
		})
	}
}
