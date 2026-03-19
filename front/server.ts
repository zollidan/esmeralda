import { existsSync, readFileSync } from "fs";
import { extname, join } from "path";

const DIST_DIR = join(import.meta.dir, "dist");
const API_URL = process.env.API_URL || "http://esmeralda-api:8080";
const PORT = Number(process.env.PORT) || 3000;

const MIME_TYPES: Record<string, string> = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".json": "application/json",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".ico": "image/x-icon",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
};

const indexHtml = readFileSync(join(DIST_DIR, "index.html"), "utf-8");

Bun.serve({
  port: PORT,
  async fetch(req) {
    const url = new URL(req.url);

    if (url.pathname.startsWith("/api")) {
      const target = `${API_URL}${url.pathname}${url.search}`;
      const headers = new Headers(req.headers);
      headers.delete("host");
      return fetch(target, {
        method: req.method,
        headers,
        body: req.body,
      });
    }

    const filePath = join(DIST_DIR, url.pathname);
    if (url.pathname !== "/" && existsSync(filePath)) {
      const ext = extname(filePath);
      return new Response(Bun.file(filePath), {
        headers: {
          "Content-Type": MIME_TYPES[ext] || "application/octet-stream",
        },
      });
    }

    return new Response(indexHtml, {
      headers: { "Content-Type": "text/html; charset=utf-8" },
    });
  },
});

console.log(`Frontend server listening on :${PORT}`);
