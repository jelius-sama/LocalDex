package api

import (
	"LocalDex/api/anime"
	"LocalDex/api/manga"
	"LocalDex/api/photo"
	m "LocalDex/middleware"
	"LocalDex/util"
	"net/http"
)

func HandleRouting() *Router {
	r := New()

	r.Route("/api/photo", func(r *Router) {
		r.Use(m.TestMiddleware, m.TestMiddlewareSecond)

		r.Post("/", photo.Post)
		r.Get("/", photo.GetMultiple)
		r.Get("/{id}", photo.Get)
		r.Delete("/", photo.DeleteMultiple)
	})

	r.Route("/api/anime", func(r *Router) {
		r.Use(m.TestMiddleware, m.TestMiddlewareSecond)

		r.Post("/", anime.Post)
		r.Get("/{id}", anime.Get)
		r.Get("/", anime.GetMultiple)
		r.Delete("/", anime.Delete)
	})

	r.Route("/api/manga", func(r *Router) {
		r.Use(m.TestMiddleware, m.TestMiddlewareSecond)

		r.Post("/", manga.Post)
		r.Get("/{id}", manga.Get)
		r.Get("/", manga.GetMultiple)
		r.Delete("/", manga.Delete)
	})

	r.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		NotFound(w, r, util.AddrOf("Route does not exists!"))
	}

	return r
}

// func parseFilter(filterStr string) map[string][]string {
//     filters := map[string][]string{}
//     if filterStr == "" {
//         return filters
//     }
//
//     parts := strings.Split(filterStr, "+")
//     for _, p := range parts {
//         kv := strings.SplitN(p, ":", 2)
//         if len(kv) == 2 {
//             key := kv[0]
//             vals := strings.Split(kv[1], ",")
//             filters[key] = append(filters[key], vals...)
//         }
//     }
//
//     return filters
// }

// filters := parseFilter("type:doujin+tags:nsfw,parody")
// filters["type"] = []string{"doujin"}
// filters["tags"] = []string{"nsfw", "parody"}
