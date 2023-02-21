package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// define routes and corresponding handlers
	http.HandleFunc("/goto", gotoHandler)
	http.HandleFunc("/return", returnHandler)

	// start the server and listen for incoming requests
	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func gotoHandler(w http.ResponseWriter, r *http.Request) {
	// get the "target" and "source" parameters from the query string
	target, err := validateURLParam(r, "target")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Bad request from %v: %v\n", r.RemoteAddr, err)
		return
	}

	source, err := validateURLParam(r, "source")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Bad request from %v: %v\n", r.RemoteAddr, err)
		return
	}

	// create a new cookie with the "source" URL
	cookie := &http.Cookie{
		Name:  "source",
		Value: source,
	}

	// set the cookie and redirect to the "target" URL
	http.SetCookie(w, cookie)
	http.Redirect(w, r, target, http.StatusSeeOther)

	fmt.Printf("Redirecting %v to %v (source=%v)\n", r.RemoteAddr, target, source)
}

func returnHandler(w http.ResponseWriter, r *http.Request) {
	// get the "source" cookie
	cookie, err := r.Cookie("source")

	// if the cookie is not set, return a 400 bad request error
	if err != nil {
		http.Error(w, "Missing source cookie", http.StatusBadRequest)
		fmt.Printf("Bad request from %v: missing source cookie\n", r.RemoteAddr)
		return
	}

	source := cookie.Value

	// clear the "source" cookie
	cookie.MaxAge = -1
	cookie.Value = ""
	http.SetCookie(w, cookie)

	// redirect to the "source" URL
	http.Redirect(w, r, source, http.StatusSeeOther)

	fmt.Printf("Redirecting %v to (source=%v)\n", r.RemoteAddr, source)
}

func validateURLParam(r *http.Request, paramName string) (string, error) {
	param := r.URL.Query().Get(paramName)
	// parse the parameter value as a URL
	u, err := url.ParseRequestURI(param)
	if err != nil {
		return "", fmt.Errorf("Invalid %v URL parameter: %v", paramName, err)
	}

	// ensure that the URL is absolute
	if !u.IsAbs() {
		return "", fmt.Errorf("Relative %v URL not allowed: %v", paramName, param)
	}

	// return the validated URL as a string
	return u.String(), nil
}
