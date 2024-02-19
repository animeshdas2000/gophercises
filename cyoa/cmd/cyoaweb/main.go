package main

import (
	"cyoa"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 3000, "The port to start CYOA applicaiton")
	fileName := flag.String("file", "gopher.json", "The JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s \n", *fileName)
	f, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", story)
	//tpl := template.Must(template.New("").Parse("Hello"))
	tpl := template.Must(template.New("").Parse(storyTemplate))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPathFunc(pathFn)) // Pass WithTemplate(tpl) to render alternative template
	mux := http.NewServeMux()
	mux.Handle("/story", h)
	fmt.Printf("Starting Port on %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTemplate = `
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
    <li><a href="/story/{{ .Chapter}}">{{ .Text}}</a></li>
    {{end}}
   </ul>

</body>

</html>
`
