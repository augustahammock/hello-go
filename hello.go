package main

import (
  "html/template"
  "log"
  "net/http"
  "os"
  "path"
  "io/ioutil"
  "encoding/json"
)

type Example struct {
  Title string
  Description string
}

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  http.HandleFunc("/", serveTemplate)  

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}

func getJSON() Example {
    bytes, err := ioutil.ReadFile("static/json/page.json")
    if err != nil {
        panic(err)
    }

    data := Example{}

    if err := json.Unmarshal(bytes, &data); err != nil {
        panic(err)
    }

    log.Println(data)
    return data
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  // w.Header().Set("Content-Type", "text/html; charset=utf-8")
  lp := path.Join("templates", "layout.html")
  fp := path.Join("templates", r.URL.Path)

  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
  if err != nil {
    if os.IsNotExist(err) {
      http.NotFound(w, r)
      return
    }
  }

  // Return a 404 if the request is for a directory
  if info.IsDir() {
    http.NotFound(w, r)
    return
  }

  tmpl, err := template.ParseFiles(lp, fp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if err := tmpl.ExecuteTemplate(w, "body", getJSON()); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}