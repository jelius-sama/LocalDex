package photo

import "net/http"

func PutMultiple(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Put Multiple Photo"))
}
