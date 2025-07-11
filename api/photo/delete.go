package photo

import "net/http"

func DeleteMultiple(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete Multiple Photo"))
}
