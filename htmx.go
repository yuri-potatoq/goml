package go_ml

/* Htmx custom Attributes */
func HxPost(url string) HTMLAttribute {
	return Attr("hx-post", DoubleQuoted, url)
}

func HxPut(url string) HTMLAttribute {
	return Attr("hx-put", DoubleQuoted, url)
}

func HxDelete(url string) HTMLAttribute {
	return Attr("hx-delete", DoubleQuoted, url)
}

func HxTarget(value string) HTMLAttribute {
	return Attr("hx-target", DoubleQuoted, value)
}

func HxSwap(value string) HTMLAttribute {
	return Attr("hx-swap", DoubleQuoted, value)
}

// make a better HTML async supporte on DSL
func HxOn(event, expr string) HTMLAttribute {
	return Attr("hx-on:"+event, DoubleQuoted, expr)
}
