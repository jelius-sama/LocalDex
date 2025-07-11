```go
/*
********************************* API Structure *****************************

Photo (supports videos as well):

 1. Add a new photo to the library:
    POST /api/photo

 2. Get a specific photo by ID:
    GET /api/photo/{id}

 3. Get multiple photos (used for homepage, browsing, etc.):
    GET /api/photo?sort={SORT}&limit={LIMIT}&page={PAGE}&filter={FILTER}

    Examples:
    - /api/photo?filter=tags:vacation+favorite:true
    - /api/photo?sort=created_desc&page=2&limit=20

 4. Delete one or more photos by ID:
    DELETE /api/photo?ids={id1,id2,id3}

 5. Recover recently deleted photos:
    PUT /api/photo/recover?ids={id1,id2}

    TODO:
    - Albums
    - Favorites
    - Captions
    - Metadata editing
    - Search

Anime (includes Hentai):

 1. Add a new anime/hentai to the library:
    POST /api/anime

 2. Get a specific anime/hentai by ID:
    GET /api/anime/{id}

 3. Get multiple anime/hentai entries:
    GET /api/anime?sort={SORT}&limit={LIMIT}&page={PAGE}&filter={FILTER}

    Examples:
    - /api/anime?filter=type:anime+tags:action+status:watching
    - /api/anime?filter=type:hentai+tags:nsfw,school&sort=added_desc
    - /api/anime?sort=title_asc&limit=15

 4. Delete a single anime/hentai by ID:
    DELETE /api/anime/{id}

    TODO:
    - Metadata editing
    - Watchlist
    - Watch history
    - Like/Dislike
    - View count
    - Custom lists
    - Search

Manga (includes Doujin):

 1. Add a new manga/doujin entry:
    POST /api/manga

 2. Get a specific manga/doujin by ID:
    GET /api/manga/{id}

 3. Get multiple entries (manga and/or doujin):
    GET /api/manga?sort={SORT}&limit={LIMIT}&page={PAGE}&filter={FILTER}

    Examples:
    - /api/manga?filter=type:manga+tags:isekai
    - /api/manga?filter=type:doujin+tags:nsfw,parody&sort=readcount_desc
    - /api/manga?sort=created_desc&limit=30&page=1

 4. Delete a single manga/doujin by ID:
    DELETE /api/manga/{id}

    TODO:
    - Metadata editing
    - Readlist
    - Read history
    - Like/Dislike
    - Read count
    - Custom lists
    - Search

*****************************************************************************
*/
```
