package manga

import "net/http"

func Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Manga"))
}

func GetMultiple(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Multiple Manga"))
}
