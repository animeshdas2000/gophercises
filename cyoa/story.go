package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTempl))
}

var defaultHandlerTempl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title}}</title>
</head>
<body> 
   <h2> {{.Title}}</h2>
   {{range .Paragraphs}}
   <p>{{.}}</p>
   {{end}}
   <ul>
    {{range .Options}}
    <li><a href="/{{ .Chapter}}">{{ .Text}}</a></li>
    {{end}}
   </ul>

</body>

</html>
`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func defaultPathFn(r *http.Request) string {

	//triming any whitespace from URL path
	path := strings.TrimSpace(r.URL.Path)

	//Default Routing to /intro
	if path == "" || path == "/" {
		path = "/intro"
	}

	//omitting "/"
	return path[1:]

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if Chapter is found and ok; Execute Chappter s
	if chapter, ok := h.s[h.pathFn(r)]; ok {
		err := h.t.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter Not Found", http.StatusNotFound)

}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
