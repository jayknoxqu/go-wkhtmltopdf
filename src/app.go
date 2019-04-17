package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func main() {
	const bindAddress = ":80"
	http.HandleFunc("/", requestHandler)
	fmt.Println("Http server listening on", bindAddress)
	http.ListenAndServe(bindAddress, nil)
}

type documentRequest struct {
	Url    string
	Urls   []string
	Output string
	// TODO: whitelist options that can be passed to avoid errors,
	// log warning when different options get passed
	Options  map[string]interface{}
	Params   []string
	Cookies  map[string]string
	FileName string
}

func logOutput(request *http.Request, message string) {
	ip := strings.Split(request.RemoteAddr, ":")[0]
	fmt.Println(ip, request.Method, request.URL, message)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		response.WriteHeader(http.StatusNotFound)
		logOutput(request, "404 not found")
		return
	}
	//if request.Method != "POST" {
	//	response.Header().Set("Allow", "POST")
	//	response.WriteHeader(http.StatusMethodNotAllowed)
	//	logOutput(request, "405 not allowed")
	//	return
	//}
	//decoder := json.NewDecoder(request.Body)
	var req documentRequest
	if err := json.Unmarshal([]byte(request.FormValue("param")), &req); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		logOutput(request, "400 bad request (invalid JSON)")
		return
	}
	segments := make([]string, 0)
	for key, element := range req.Options {
		if element == true {
			// if it was parsed from the JSON as an actual boolean,
			// convert to command-line single argument	(--foo)
			segments = append(segments, fmt.Sprintf("--%v", key))
		} else if element != false {
			// Otherwise, use command-line argument with value (--foo bar)
			segments = append(segments, fmt.Sprintf("--%v", key), fmt.Sprintf("%v", element))
		}
	}

	if len(req.Params) > 0 {
		for _, param := range req.Params {
			segments = append(segments, param)
		}
	}

	for key, value := range req.Cookies {
		segments = append(segments, "--cookie", key, url.QueryEscape(value))
	}
	var programFile string
	var contentType string
	var fileSuffix string
	switch req.Output {
	case "jpg":
		programFile = "/bin/wkhtmltoimage"
		contentType = "image/jpeg"
		fileSuffix = ".jpeg"
		segments = append(segments, "--format", "jpg", "-q")
	case "png":
		programFile = "/bin/wkhtmltoimage"
		contentType = "image/png"
		fileSuffix = ".png"
		segments = append(segments, "--format", "png", "-q")
	default:
		// defaults to pdf
		programFile = "/bin/wkhtmltopdf"
		contentType = "application/pdf"
		fileSuffix = ".pdf"
	}
	if req.Url != "" {
		u, _ := url.QueryUnescape(req.Url)
		segments = append(segments, u, "-")
	} else if len(req.Urls) > 0 {
		for _, u := range req.Urls {
			uu, _ := url.QueryUnescape(u)
			segments = append(segments, uu)
		}
		segments = append(segments, "-")
	}
	fmt.Println("\tRunning:", programFile, strings.Join(segments, " "))
	cmd := exec.Command(programFile, segments...)
	response.Header().Set("Content-Type", contentType)
	if req.FileName != "" {
		response.Header().Set("Content-Disposition", "attachment;filename="+url.QueryEscape(req.FileName)+fileSuffix)
	}
	cmd.Stdout = response
	cmd.Start()
	defer cmd.Wait()
	// TODO: check if Stderr has anything, and issue http 500 instead.
	logOutput(request, "200 OK")
}
