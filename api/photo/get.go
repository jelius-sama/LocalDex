package photo

import "net/http"

func Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Photo"))
}

func GetMultiple(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Multiple Photo"))
}
