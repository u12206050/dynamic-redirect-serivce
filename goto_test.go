package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGotoHandler(t *testing.T) {
	// create a new request with the "target" and "source" parameters set
	req, err := http.NewRequest("GET", "/goto?target=http://example.com&source=http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// call the handler function with the request and recorder
	handler := http.HandlerFunc(gotoHandler)
	handler.ServeHTTP(rr, req)

	// ensure that the response status code is correct
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusSeeOther)
	}

	// ensure that the cookie was set with the correct name and value
	cookie := rr.Header().Get("Set-Cookie")
	if cookie != "source=http://localhost" {
		t.Errorf("Handler set wrong cookie: got %v, want %v", cookie, "source=http://localhost")
	}

	// ensure that the response was redirected to the correct URL
	location := rr.Header().Get("Location")
	if location != "http://example.com" {
		t.Errorf("Handler redirected to wrong URL: got %v, want %v", location, "http://example.com")
	}
}

func TestReturnHandler(t *testing.T) {
	// create a new request with no cookies set
	req, err := http.NewRequest("GET", "/return", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// call the handler function with the request and recorder
	handler := http.HandlerFunc(returnHandler)
	handler.ServeHTTP(rr, req)

	// ensure that the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusBadRequest)
	}

	// create a new request with a valid "source" cookie set
	cookie := &http.Cookie{Name: "source", Value: "http://localhost"}
	req, err = http.NewRequest("GET", "/return", nil)
	req.AddCookie(cookie)

	// create a new recorder to capture the response
	rr = httptest.NewRecorder()

	// call the handler function with the request and recorder
	handler = http.HandlerFunc(returnHandler)
	handler.ServeHTTP(rr, req)

	// ensure that the response status code is correct
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusSeeOther)
	}

	// ensure that the cookie was cleared
	cookies := rr.Header().Values("Set-Cookie")
	if len(cookies) != 1 || cookies[0] != "source=; Max-Age=0" {
		t.Errorf("Handler didn't clear cookie: got %v, want %v", cookies, []string{"source=; Max-Age=0"})
	}

	// ensure that the response was redirected to the correct URL
	location := rr.Header().Get("Location")
	if location != "http://localhost" {
		t.Errorf("Handler redirected to wrong URL: got %v, want %v", location, "http://localhost")
	}
}
