package main 

import (
	"io"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"html/template"
	"path"
	"runtime/debug"
)

const (
	UPLOAD_DIR = "./uploads"
	VIEWS_DIR = "./views"
)

var templates map[string]*template.Template = make(map[string]*template.Template)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		/*io.WriteString(w, "<html><form method=\"POST\" action=\"/upload\" " + 
			"enctype=\"multipart/form-data\">" + 
			"Choose an image to upload: <input name=\"image\" type=\"file\" />" + 
			"<input type=\"submit\" value=\"Upload\" />" + 
			"</form></html>")*/
		//t, err := template.ParseFiles(VIEWS_DIR + "/upload.html")
		if err := renderHtml(w, "upload.html", nil); err != nil {
			check(err)	
		}
		//t.Execute(w, nil)	
		return
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			check(err)	
		}

		filename := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			check(err)	
		}

		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			check(err)
		}

		http.Redirect(w, r, "/view?id=" + filename, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	} 
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir(UPLOAD_DIR)
	if err != nil {
		check(err)	
	}

	locals := make(map[string]interface{})
	images := []string{}

	//var listHtml string
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
		//imgid := fileInfo.Name()
		//listHtml += "<li><a href=\"/view?id=" + imgid + "\">" + imgid  + "</a></li>"
	}
	locals["images"] = images
	//t, err := template.ParseFiles(VIEWS_DIR + "/list.html")
	if err := renderHtml(w, "list.html", locals); err != nil {
		check(err)	
	}
	// t.Execute(w, locals)

	//io.WriteString(w, "<html><ol>" + listHtml + "</ol></html>")
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50X页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error.html", e)
				// logging
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}	
		}()
		fn(w, r)
	}
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) (err error) {
	// t, err := template.ParseFiles(VIEWS_DIR + "/" + tmpl + ".html")
	// if err != nil {
	// 	return
	// }
	// err = t.Execute(w, locals)
	err = templates[tmpl].Execute(w, locals)
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	fileInfoArr, err := ioutil.ReadDir(VIEWS_DIR)
	if err != nil {
		panic(err)
		return
	}

	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue	
		}
		templatePath = VIEWS_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}
}

func main() {
	http.HandleFunc("/", safeHandler(listHandler))
	http.HandleFunc("/upload", safeHandler(uploadHandler))
	http.HandleFunc("/view", safeHandler(viewHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}