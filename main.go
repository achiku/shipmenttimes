package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func sendError(w http.ResponseWriter, e error, code int) {
	fmt.Fprintf(w, http.StatusText(code))
	w.WriteHeader(code)
	log.Printf("error: %s", e)
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	lns, err := ParseBaseCSV(file)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	odrs, err := TransformBaseOrder(lns)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	dt := time.Now().Format("20060102")
	basepath := filepath.Join("output", dt)
	if err := os.MkdirAll(basepath, 0744); err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	clickpost, other := QuantityFilter(odrs, 4)

	cf, err := os.Create(filepath.Join(basepath, "clickpost.csv"))
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	if err := WriteClickpostFormat(cf, clickpost); err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	of, err := os.Create(filepath.Join(basepath, "other.csv"))
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	if err := WriteSummaryFormat(of, other); err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	sf, err := os.Create(filepath.Join(basepath, "summary.csv"))
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	if err := WriteSummaryFormat(sf, odrs); err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `
<html>
  <header>
    <title>SHIPMENT-TIMES</title>
  </header>
  <body>
    <div>
      <form enctype="multipart/form-data" action="http://127.0.0.1:8080/upload" method="post">
        <input type="file" name="uploadfile" />
        <input type="submit" value="upload" />
      </form>
    </div>
  </body>
</html>
	`
	fmt.Fprint(w, html)
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/upload", upload)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
