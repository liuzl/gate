package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func a(w http.ResponseWriter, r *http.Request) {
	showError(w, r, "error", 400)
	return
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Fprintln(w, fmt.Sprintf("%v=%v", k, v))
	}
	fmt.Fprintln(w, r.Form.Get("xyz"))

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintln(w, string(b))
	}
	mustEncode(w, r.Form)
}

func main() {
	http.Handle("/a", http.HandlerFunc(a))
	http.ListenAndServe(":8383", nil)
}
