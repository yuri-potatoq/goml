package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ht "github.com/yuri-potatoq/go_ml"
)

var (
	htmxCDN     = "https://unpkg.com/htmx.org@1.9.9"
	tailwindCDN = "https://cdn.tailwindcss.com"
)

/* Data */
type TodoList struct {
	id        int
	title     string
	isChecked bool
}

type TodoListDb struct {
	storage []TodoList
}

func (db *TodoListDb) Add(td TodoList) []TodoList {
	td.id = len(db.storage)

	fmt.Printf("Before Add: %v\n", db.storage)
	db.storage = append(db.storage, td)
	fmt.Printf("After Add: %v\n", db.storage)
	return db.storage
}

func (db *TodoListDb) Update(todoId int, isChecked bool) error {
	fmt.Printf("Before Update: %v\n", db.storage)
	if todoId < 0 || todoId > len(db.storage) {
		return fmt.Errorf("not found with values: Id: %d Checked: %v", todoId, isChecked)
	}

	current := db.storage[todoId]
	current.isChecked = isChecked
	db.storage[todoId] = current

	fmt.Printf("After Update: %v\n", db.storage)
	return nil
}

func (db *TodoListDb) Get(todoId int) (TodoList, error) {
	if todoId >= 0 && todoId <= len(db.storage) {
		return db.storage[todoId], nil
	}
	return TodoList{}, errors.New("not found")
}

func (db *TodoListDb) GetAll() []TodoList {
	return db.storage
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
					ht.Id(fmt.Sprintf("%d", t.id)),
					ht.IsChecked(t.isChecked),
					ht.HxOn("click", "fetch(`/todo/${this.checked ? 'enable' : 'disable'}/${this.id}`, {method: 'PUT'})")),
			),
			ht.Th()(ht.RawText(t.title)),
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
	db := new(TodoListDb)
	atoiOrNegative := func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			return -1
		}
		return n
	}

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
			fmt.Println(pathStrs)
			// should only be /todo/{disable|enable}/{n}
			if err := db.Update(atoiOrNegative(pathStrs[3]), pathStrs[2] == "enable"); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			break
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
