package anime

import "net/http"

func Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Anime"))
}

func GetMultiple(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Multiple Anime"))
}
