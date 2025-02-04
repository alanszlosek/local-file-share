package main

import (
    _ "embed"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)


//go:embed "index.tmpl"
var indexTemplate string

var destination string = "uploads"

// https://dev.to/neelp03/building-a-file-upload-service-in-go-34fj

func main() {
    http.HandleFunc("/upload", fileUploadHandler)
    http.HandleFunc("/", defaultHandler)

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

type Data struct {
    Files []string
}
func defaultHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
    path := httpRequest.URL.Path

    if strings.HasPrefix(path, "/get/") {
        fileHandler(httpResponse, httpRequest)
        return
    } 

    homeHandler(httpResponse, httpRequest)

}

func homeHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
    // Read files from destination dir, and send to template
    entries, err := os.ReadDir(destination)
    if err != nil {
        log.Fatal(err)
    }
    files := []string{}
    for _, entry := range entries {
        files = append(files, entry.Name())
    }
    fmt.Printf("%#v\n", files)

    data := Data{
        Files: files,
    }
    tmpl, err := template.New("template").Parse(indexTemplate)
	if err != nil {
		fmt.Printf("%s\n", err)
        http.Error(httpResponse, "Error loading template", http.StatusInternalServerError)
		return
	}
	var rendered strings.Builder
	err = tmpl.Execute(&rendered, data)
	if err != nil {
		fmt.Printf("%s\n", err)
		http.Error(httpResponse, "Error rendering template", http.StatusInternalServerError)
        return
	}

    httpResponse.Header().Add("Content-type", "text/html")
    httpResponse.Write( []byte(rendered.String()) )
}

func fileHandler(httpResponse http.ResponseWriter, r *http.Request) {
    filename := strings.ReplaceAll(r.URL.Path, "/get/", "")
    fmt.Printf("Reading %s\n", filename)

    contents, err := os.ReadFile(filepath.Join(destination, filename))
    if err != nil {
        fmt.Printf("%s\n", err)
		http.Error(httpResponse, "Error rendering template", http.StatusInternalServerError)
        return
    }

    mimeType := http.DetectContentType(contents)

    httpResponse.Header().Add("Content-type", mimeType)
    httpResponse.Write( contents )

}

func fileUploadHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
    // Limit file size to 10MB. This line saves you from those accidental 100MB uploads!
    httpRequest.ParseMultipartForm(10 << 20)

    // Retrieve the file from form data
    file, handler, err := httpRequest.FormFile("myfile")
    if err != nil {
        fmt.Printf("Failed to process upload: %s\n", err)
        http.Error(httpResponse, "Error retrieving the file: %s\n", http.StatusBadRequest)
        return
    }
    defer file.Close()



    // Now let’s save it locally
    filename := handler.Filename
    for _,p := range []string{"/", "\\\\", ".."} {
        filename = strings.ReplaceAll(filename, p, "")
    }
    dst, err := createFile(filename)
    if err != nil {
        http.Error(httpResponse, "Error saving the file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy the uploaded file to the destination file
    if _, err := dst.ReadFrom(file); err != nil {
        http.Error(httpResponse, "Error saving the file", http.StatusInternalServerError)
    }

    out := fmt.Sprintf(
        `Uploaded File: %s<br />
File Size: %d<br />
MIME Header: %v<br /><br />
<a href="/">Home</a>
`,
        handler.Filename,
        handler.Size,
        handler.Header,
    )
    httpResponse.Header().Add("Content-type", "text/html")
    httpResponse.Write( []byte(out) )
}

func createFile(filename string) (*os.File, error) {
    // Create an uploads directory if it doesn’t exist
    if _, err := os.Stat(destination); os.IsNotExist(err) {
        os.Mkdir(destination, 0755)
    }

    // Build the file path and create it
    dst, err := os.Create(filepath.Join(destination, filename))
    if err != nil {
        return nil, err
    }

    return dst, nil
}


