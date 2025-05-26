package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var zoneID = os.Getenv("CLOUDFLARE_ZONEID")
var pat = os.Getenv("CLOUDFLARE_PAT")

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	http.HandleFunc("/api/CloudflareCachebuster", purgeCacheHandler)
	log.Printf("Listening on %s. Serving at http://127.0.0.1:%s/", customHandlerPort, customHandlerPort)

	log.Fatal(http.ListenAndServe(":"+customHandlerPort, nil))
}

// httpError logs the error and returns an HTTP error message and code.
func httpError(w http.ResponseWriter, err error, msg string, errorCode int) {
	errorMsg := fmt.Sprintf("%s: %v", msg, err)
	log.Printf("%s", errorMsg)
	http.Error(w, errorMsg, errorCode)
}

func purgeCacheHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received %s from %v", r.Method, r.RemoteAddr)
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpError(w, err, "error reading POST body", http.StatusInternalServerError)
			return
		}
		log.Printf("Request body: %s", body)
	}
	// Send POST request to Cloudflare
	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zoneID)
	data := `{"purge_everything":true}`
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data))
	if err != nil {
		httpError(w, err, "error creating new Request", http.StatusInternalServerError)
		return
	}

	authHeader := fmt.Sprintf("Bearer %s", pat)
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	cloudflareResp, err := client.Do(req)
	if err != nil {
		httpError(w, err, "error sending POST request", http.StatusInternalServerError)
		return
	}
	defer func() {
		err := cloudflareResp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Pass cloudflare response to caller

	cloudflareRespBody, err := io.ReadAll(cloudflareResp.Body)
	if err != nil {
		httpError(w, err, "error reading Cloudflare response", http.StatusInternalServerError)
		return
	}

	if cloudflareResp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("error non-200 status: %s", cloudflareRespBody)
		httpError(w, nil, msg, http.StatusInternalServerError)
		return
	}

	log.Printf("Cloudflare response: %s", cloudflareRespBody)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(cloudflareRespBody)
	if err != nil {
		httpError(w, err, "error sending response to client", http.StatusInternalServerError)
		return
	}
}
