package main

import (
	"net/http"

	"github.com/yuri-potatoq/go_ml"
)

func HelloWorldView() go_ml.HTMLContent {
	return go_ml.Html(go_ml.Lang("en"))(
		go_ml.Div(go_ml.ClassNames("container"))(
			go_ml.RawText("Hello World!"),
		),
	)
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(HelloWorldView().BuildDOM()))
	}))
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
