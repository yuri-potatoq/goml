package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ht "github.com/yuri-potatoq/go_ml"
)

var (
	htmxCDN     = "https://unpkg.com/htmx.org@1.9.9"
	tailwindCDN = "https://cdn.tailwindcss.com"
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
func PageIndex(partials ...ht.HTMLContent) ht.HTMLContent {
	return ht.Html()(
		ht.Head()(
			ht.Title()(ht.RawText("Todo List")),
			// tag with custom attributes
			ht.Script(ht.Src(htmxCDN), ht.Attr("crossorigin", "anonymous"))(),
			ht.Script(ht.Src(tailwindCDN))(),
		),
		ht.Body(
			ht.ClassNames("h-100 w-full flex items-center justify-center bg-teal-lightest font-sans"),
		)(ht.Div()(partials...)),
	)
}

func ListOfTodos(todos ...TodoList) ht.HTMLContent {
	if len(todos) == 0 {
		return ht.Div(ht.Id("todo-list-tb-container"))()
	}

	var todoLines []ht.HTMLContent
	for _, t := range todos {
		todoLines = append(todoLines, ht.Tr()(
			ht.Th(ht.ClassNames("h-10 border border-gray-100 rounded"))(
				ht.Input(
					ht.Type("checkbox"),
					ht.Id(t.id),
					ht.IsChecked(t.isChecked),
					ht.HxOn("click", "fetch(`/todo/${this.checked ? 'enable' : 'disable'}/${this.id}`, {method: 'PUT'})")),
			),
			ht.Th()(ht.RawText(t.title)),
			ht.Th()(
				ht.Button(
					ht.ClassNames("bg-white hover:bg-gray-100 text-gray-800 font-semibold py-2 px-4 border border-gray-400 rounded hover:shadow-inner"),
					ht.HxDelete(fmt.Sprintf("/todo/%s", t.id)),
					ht.HxTarget("#todo-list-tb-container"),
					ht.HxSwap("outerHTML"),
				)(ht.RawText("Delete")),
			),
		))
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
				ht.ClassNames("shadow appearance-none border rounded w-full py-2 px-3 mr-4 text-grey-darker")),
			ht.Div(ht.ClassNames("flex w-full items-center justify-center p-1"))(
				ht.Button(
					ht.ClassNames("p-2 border-2 rounded text-teal border-teal hover:text-white hover:bg-teal"),
					ht.Type("submit"))(ht.RawText("Submit")),
			),
		),
	)
}

func main() {
	db := TodoListDb{storage: make(map[string]TodoList)}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathStrs := strings.Split(strings.TrimSpace(r.URL.Path), "/")

		switch r.Method {
		case http.MethodGet:
			doc := PageIndex(AddTodoForm(), ListOfTodos(db.GetAll()...))
			w.Write([]byte(doc.BuildDOM()))
			break
		case http.MethodPost:
			doc := ListOfTodos(db.Add(TodoList{
				title:     r.FormValue("task"),
				isChecked: false,
			})...)
			w.Write([]byte(doc.BuildDOM()))
			break
		case http.MethodPut:
			// should only be /todo/{disable|enable}/{n}
			if err := db.Update(pathStrs[3], pathStrs[2] == "enable"); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		case http.MethodDelete:
			db.Delete(pathStrs[2])
			w.Write([]byte(ListOfTodos(db.GetAll()...).BuildDOM()))
			break
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
