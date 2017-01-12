package main

import (
  "cloud.google.com/go/storage"
  "encoding/json"
  "golang.org/x/net/context"
  "google.golang.org/appengine"
  "html/template"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "path"
)

// var globalLayout `{{define "page"}}

// <!doctype html>
// <html>
// {{template "head" .Data}}

// <body>
//   {{template "global-nav" .}}

//   {{template `  `}}

//   {{template "footer" .}}

// </body>
// </html>
// {{end}}`

type GlobalData struct {
}

// type PageData struct {
//   GlobalData   GlobalData
//   CssResources []string
//   JsResources  []string
//   Content      map[string]interface{}
// }

type PageData struct {
  GlobalData GlobalData
  Content    map[string]interface{}
}

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  http.HandleFunc("/", serveTemplate)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
  appengine.Main()
}

func getJSON() map[string]interface{} {

  ctx := context.Background()
  client, err := storage.NewClient(ctx)
  if err != nil {
    log.Println("error creating client")
  } else {
    log.Println("created client")
  }

  bkt := client.Bucket("augustas-static-assets")
  if err := bkt.Create(ctx, "storage-project-151921", nil); err != nil {
    log.Println("error connecting to bucket")
  } else {
    log.Println("connected to bucket")
  }

  attrs, err := bkt.Attrs(ctx)
  if err != nil {
    log.Println("error with bucket attributes")
  } else {
    log.Println("bucket attributes OK")
  }
  log.Println("bucket %s, created at %s, is located in %s with storage class %s\n",
    attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)

  obj := bkt.Object("content/development/evernote-sandbox/page.json")

  r, err := obj.NewReader(ctx)
  if err != nil {
    log.Println("could not read file")
  } else {
    log.Println("read file OK")
    // log.Println(r)
  }
  defer r.Close()

  var jsonData = make(map[string]interface{})

  slurp, err := ioutil.ReadAll(r)
  if err != nil {
    log.Println(err)
  }

  if err := json.Unmarshal(slurp, &jsonData); err != nil {
    log.Println("test")
    panic(err)
  } else {
    log.Println(jsonData)
  }

  // bytes, err := ioutil.ReadFile("static/json/page.json")
  // if err != nil {
  //   panic(err)
  // }

  // data := Example{}

  // if err := json.Unmarshal(bytes, &data); err != nil {
  //   panic(err)
  // }

  // log.Println(data)
  return jsonData
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  getJSON()

  // w.Header().Set("Content-Type", "text/html; charset=utf-8")
  // tp := path.Join("templates", "top.html")
  // hp := path.Join("templates", "header.html")
  // fp := path.Join("templates", "footer.html")
  fp := path.Join("templates", r.URL.Path)
  log.Println(r.URL.Path)

  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
  if err != nil {
    if os.IsNotExist(err) {
      panic("Template not found")
      http.NotFound(w, r)
      return
    }
  }

  // Return a 404 if the request is for a directory
  if info.IsDir() {
    panic("Directory not found")
    http.NotFound(w, r)
    return
  }

  // tmpl, err := template.ParseFiles(fp)
  tmpl := template.Must(template.ParseGlob("templates/*.tmpl"))

  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  jsonData := getJSON()
  pageData := PageData{}
  pageData.Content = jsonData

  if err := tmpl.ExecuteTemplate(w, "page", pageData); err != nil {
    // if err := tmpl.ExecuteTemplate(w, "page", getJSON()); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}
