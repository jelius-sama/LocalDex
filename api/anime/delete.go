package anime

import "net/http"

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete Anime"))
}
