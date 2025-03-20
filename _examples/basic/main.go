package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/dreams-money/turnstile"
)

var (
	siteVerifySecretCode = os.Getenv("TURNSTILE-SECRET")
	siteVerifyClientCode = os.Getenv("TURNSTILE-PUBLIC")
)

var (
	botcheck = turnstile.New(siteVerifySecretCode)
	login    = template.Must(template.ParseFiles("index.html"))
)

func index(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	login.Execute(resp, map[string]string{
		"TurnstileSiteKey": siteVerifyClientCode,
	})
}

func submit(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(resp, "Server Error", http.StatusInternalServerError)
		return
	}

	botErr, requestErr := botcheck.Verify(req.FormValue("cf-turnstile-response"), req.RemoteAddr)
	if requestErr != nil {
		log.Println(requestErr)
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	if botErr == turnstile.ErrTimeoutOrDuplicate {
		http.Error(resp, "Refresh form before resubmitting it", http.StatusBadRequest)
		return
	} else if botErr != nil { // Turnstile found a bot
		log.Println("Bot match: ", err)
		http.Error(resp, "Are you a bot?", http.StatusUnauthorized)
		return
	}

	name := req.FormValue("name")

	fmt.Fprintf(resp, "Your name is: %v", name)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/submit", submit)

	log.Println("Server starting at http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}
