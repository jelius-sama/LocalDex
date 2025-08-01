const CACHE_NAME = "localdex-pwa-cache-v1";

const urlsToCache = [
  "/",
  "/assets/favicon.png",
  "/assets/manifest.json",
  "/assets/icons/apple-icon-180.png",
  "/assets/icons/apple-splash-2048-2732.png",
  "/assets/icons/apple-splash-2732-2048.png",
  "/assets/icons/apple-splash-1668-2388.png",
  "/assets/icons/apple-splash-2388-1668.png",
  "/assets/icons/apple-splash-1536-2048.png",
  "/assets/icons/apple-splash-2048-1536.png",
  "/assets/icons/apple-splash-1488-2266.png",
  "/assets/icons/apple-splash-2266-1488.png",
  "/assets/icons/apple-splash-1640-2360.png",
  "/assets/icons/apple-splash-2360-1640.png",
  "/assets/icons/apple-splash-1668-2224.png",
  "/assets/icons/apple-splash-2224-1668.png",
  "/assets/icons/apple-splash-1620-2160.png",
  "/assets/icons/apple-splash-2160-1620.png",
  "/assets/icons/apple-splash-1290-2796.png",
  "/assets/icons/apple-splash-2796-1290.png",
  "/assets/icons/apple-splash-1179-2556.png",
  "/assets/icons/apple-splash-2556-1179.png",
  "/assets/icons/apple-splash-1284-2778.png",
  "/assets/icons/apple-splash-2778-1284.png",
  "/assets/icons/apple-splash-1170-2532.png",
  "/assets/icons/apple-splash-2532-1170.png",
  "/assets/icons/apple-splash-1125-2436.png",
  "/assets/icons/apple-splash-2436-1125.png",
  "/assets/icons/apple-splash-1242-2688.png",
  "/assets/icons/apple-splash-2688-1242.png",
  "/assets/icons/apple-splash-828-1792.png",
  "/assets/icons/apple-splash-1792-828.png",
  "/assets/icons/apple-splash-1242-2208.png",
  "/assets/icons/apple-splash-2208-1242.png",
  "/assets/icons/apple-splash-750-1334.png",
  "/assets/icons/apple-splash-1334-750.png",
  "/assets/icons/apple-splash-640-1136.png",
  "/assets/icons/apple-splash-1136-640.png",
];

// Install event — pre-cache important files
self.addEventListener("install", (event) => {
  self.skipWaiting(); // Force activation after install
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
      .catch((err) => console.error("Cache install failed", err)),
  );
});

// Activate event — clean up old caches
self.addEventListener("activate", (event) => {
  clients.claim(); // Take control of pages ASAP
  event.waitUntil(
    caches.keys().then((cacheNames) =>
      Promise.all(
        cacheNames.map((name) => {
          if (name !== CACHE_NAME) {
            return caches.delete(name);
          }
        }),
      ),
    ),
  );
});

// Fetch event — try cache first, fall back to network
self.addEventListener("fetch", (event) => {
  if (event.request.method !== "GET") return;

  event.respondWith(
    caches.match(event.request).then((cached) => {
      return (
        cached ||
        fetch(event.request)
          .then((response) => {
            const cloned = response.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(event.request, cloned);
            });
            return response;
          })
          .catch(() => {
            // Optional: return a fallback offline page here
          })
      );
    }),
  );
});
