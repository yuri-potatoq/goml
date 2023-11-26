package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	ht "github.com/yuri-potatoq/go_ml"
)

var (
	htmxCDN     = "https://unpkg.com/htmx.org@1.9.9"
	tailwindCDN = "https://cdn.tailwindcss.com"

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	db = TodoListDb{storage: make(map[string]TodoList)}
)

/* Data */
type TodoList struct {
	id        string
	title     string
	isChecked bool
}

type TodoListDb struct {
	storage map[string]TodoList
}

func (db *TodoListDb) Add(td TodoList) (tds []TodoList) {
	td.id = fmt.Sprintf("%d", time.Now().Unix())
	db.storage[td.id] = td
	for _, t := range db.storage {
		tds = append(tds, t)
	}
	return
}

func (db *TodoListDb) Delete(todoId string) {
	delete(db.storage, todoId)
}

func (db *TodoListDb) Update(todoId string, isChecked bool) error {
	current, ok := db.storage[todoId]
	if !ok {
		return fmt.Errorf("[%s] not found!", todoId)
	}
	current.isChecked = isChecked
	db.storage[todoId] = current
	return nil
}

func (db *TodoListDb) EditTitle(todoId string, newTitle string) (TodoList, error) {
	current, ok := db.storage[todoId]
	if !ok {
		return TodoList{}, fmt.Errorf("[%s] not found!", todoId)
	}
	current.title = newTitle
	db.storage[todoId] = current
	return current, nil
}

func (db *TodoListDb) Get(todoId string) (TodoList, error) {
	current, ok := db.storage[todoId]
	if !ok {
		return TodoList{}, fmt.Errorf("[%s] not found!", todoId)
	}
	return current, nil
}

func (db *TodoListDb) GetAll() (tds []TodoList) {
	for _, t := range db.storage {
		tds = append(tds, t)
	}
	return
}

/* Views */
var (
	// We can declared global shared styles and reuse them (or just used to group styles) :)
	FlexContainerFull = ht.ClassNames("flex items-center justify-center w-full")
	DefaultBorder     = ht.ClassNames("border rounded")
	TodoRowButton     = ht.ClassNames("bg-white hover:bg-gray-100 text-gray-800 font-semibold py-2 px-4 border-gray-400 hover:shadow-inner")
)

// Easy component builders
func NewButton(innerText string, attrs ...ht.HTMLAttribute) ht.HTMLContent {
	return ht.Button(
		append(attrs, DefaultBorder, TodoRowButton)...,
	)(ht.RawText(innerText))
}

func PageIndex(partials ...ht.HTMLContent) ht.HTMLContent {
	return ht.Html()(
		ht.Head()(
			ht.Title()(ht.RawText("Todo List")),
			// tag with custom attributes
			ht.Script(ht.Src(htmxCDN), ht.Attr("crossorigin", "anonymous"))(),
			ht.Script(ht.Src(tailwindCDN))(),
		),
		ht.Body(
			FlexContainerFull,
			ht.ClassNames("h-100 bg-teal-lightest font-sans"),
		)(ht.Div()(partials...)),
	)
}

func EditTodoRow(t TodoList) ht.HTMLContent {
	rowId := "todo-row-" + t.id
	reqBody := fmt.Sprintf("body: `title=${document.querySelectorAll('#%s > th > input')[0].value}`", rowId)
	return ht.Tr(ht.Id(rowId), ht.ClassNames("flex justify-stretch"))(
		ht.Th()(
			ht.Input(ht.Type("text"), ht.Value(t.title)),
		),
		ht.Th()(
			NewButton("Ok",
				ht.HxOn("click",
					fmt.Sprintf("fetch(`/todo/edit/%s`, { %s, %s, %s })",
						t.id, "method: 'PUT'",
						"headers: { 'Content-Type': 'application/x-www-form-urlencoded'}",
						reqBody)),
				ht.HxGet("/todo/"+t.id),
				ht.HxTrigger("click delay:0.5s"),
				ht.HxTarget("#"+rowId),
				ht.HxSwap("outerHTML")),
		),
	)
}

func LoadTodoRow(t TodoList) ht.HTMLContent {
	rowName := "todo-row-" + t.id

	return ht.Tr(ht.Id(rowName), ht.ClassNames("flex justify-stretch"))(
		ht.Th(DefaultBorder, ht.ClassNames("h-10 border-gray-100"))(
			ht.Input(
				ht.Type("checkbox"),
				ht.Id(t.id),
				ht.IsChecked(t.isChecked),
				ht.HxOn("click",
					"fetch(`/todo/${this.checked ? 'enable' : 'disable'}/${this.id}`, {method: 'PUT'})")),
		),
		ht.Th()(ht.RawText(t.title)),
		ht.Th()(
			NewButton("Edit",
				ht.HxGet("/todo/edit/"+t.id),
				ht.HxTarget("#"+rowName),
				ht.HxSwap("outerHTML")),
			NewButton("Delete",
				ht.HxDelete("/todo/edit/"+t.id),
				ht.HxTarget("#todo-list-tb-container"),
				ht.HxSwap("outerHTML")),
		),
	)
}

func ListOfTodos(todos ...TodoList) ht.HTMLContent {
	if len(todos) == 0 {
		return ht.Div(ht.Id("todo-list-tb-container"))()
	}

	var todoLines []ht.HTMLContent
	for _, t := range todos {
		todoLines = append(todoLines,
			ht.Div(ht.ClassNames("w-full"))(LoadTodoRow(t)))
	}

	return ht.Div(ht.Id("todo-list-tb-container"))(
		ht.Table(ht.ClassNames("w-full whitespace-nowrap"))(todoLines...),
	)
}

func AddTodoForm() ht.HTMLContent {
	return ht.Form(
		ht.HxPost("/todo"),
		ht.HxTarget("#todo-list-tb-container"),
		ht.HxSwap("outerHTML"),
	)(
		ht.Div(ht.ClassNames("flex-column"))(
			ht.Input(
				ht.Name("task"),
				ht.PlaceHolder("Your task name"),
				DefaultBorder,
				ht.ClassNames("shadow appearance-none w-full py-2 px-3 mr-4 text-grey-darker")),
			ht.Div(FlexContainerFull, ht.ClassNames("p-1"))(
				NewButton("Submit",
					ht.ClassNames("p-2 text-teal border-teal hover:text-white hover:bg-teal"),
					ht.Type("submit")),
			),
		),
	)
}

/* Handlers */
func mainHandler(w http.ResponseWriter, r *http.Request) {
	buildWithOpts := func(doc ht.HTMLContent) {
		err := doc.BuildDOM(
			ht.WithDefaultIndentation(),
			ht.WithWriter(w),
			ht.WithLogger(logger),
		)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{ "message": "%s"}`, err)))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	pathStrs := strings.Split(strings.TrimSpace(r.URL.Path), "/")

	switch r.Method {
	case http.MethodGet:
		var doc ht.HTMLContent
		if len(pathStrs) >= 4 && pathStrs[2] == "edit" {
			// /todo/edit/{n}
			t, err := db.Get(pathStrs[3])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			doc = EditTodoRow(t)
		} else if len(pathStrs) >= 3 {
			// /todo/{n}
			t, err := db.Get(pathStrs[2])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			doc = LoadTodoRow(t)
		} else {
			doc = PageIndex(
				AddTodoForm(), ListOfTodos(db.GetAll()...))
		}

		buildWithOpts(doc)
		break
	case http.MethodPost:
		buildWithOpts(ListOfTodos(db.Add(TodoList{
			title:     r.FormValue("task"),
			isChecked: false,
		})...))
		break
	case http.MethodPut:
		if len(pathStrs) < 4 {
			break
		}
		action := pathStrs[2]

		switch action {
		case "edit":
			// /todo/edit/{n}
			curr, err := db.EditTitle(pathStrs[3], r.FormValue("title"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}
			buildWithOpts(EditTodoRow(curr))

			break
		case "disable", "enable":
			// should only be /todo/{disable|enable}/{n}
			if err := db.Update(pathStrs[3], action == "enable"); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		}
		break
	case http.MethodDelete:
		db.Delete(pathStrs[2])
		buildWithOpts(ListOfTodos(db.GetAll()...))
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	http.Handle("/", http.HandlerFunc(mainHandler))

	http.ListenAndServe(":8080", http.DefaultServeMux)
}
