package main

import (
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
			name:         "build script tag with non key-valued attribute",
			expectedHtml: `<script defer></script>`,
			givenDOM:     Script(Defer())(),
		},
		{
			name:         "build input regular tag",
			expectedHtml: `<input type="text"/>`,
			givenDOM:     Input(Type("text")),
		},
		{
			name:         "build non-void tag with inner text - 1",
			expectedHtml: `<div>olá &#10; mundo</div>`,
			givenDOM:     Div()(RawText("olá &#10; mundo")),
		},
		{
			name:         "build non-void tag with inner text - 2",
			expectedHtml: `<div>olá &#10; mundo<div></div>test<div></div></div>`,
			givenDOM:     Div()(RawText("olá &#10; mundo"), Div()(), RawText("test"), Div()()),
		},
		// TODO: fix wrong attribute spaces sort
		// i.g.: <input type="checkbox"required required="required"/>
		// {
		// 	name:         "build all boolean attributes possibilities",
		// 	expectedHtml: `<input type="text" required required="required"/>`,
		// 	givenDOM:     Input(Type("checkbox"), Required(), attr("required", "required")),
		// },
	}

	for _, tc := range testSuite {
		t.Run(tc.name, func(t *testing.T) {
			parsedDoc := tc.givenDOM.BuildDOM()
			if strings.Compare(parsedDoc, tc.expectedHtml) != 0 {
				t.Errorf("result not match: given: [%s], expected: [%s]", parsedDoc, tc.expectedHtml)
			}
		})
	}
}
