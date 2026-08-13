# PandalilySite

Fan site for pandabearlily, replacing her current [carrd site](https://pandabearlily.carrd.co/) (bad on mobile, hard to read, missing features).

Private repo — not open to contributions.

## Stack

Monorepo with frontend and backend kept fully separate:

- `web/` — [SvelteKit](https://kit.svelte.dev/), TypeScript, plain CSS (no Tailwind). Most pages prerendered/static.
- `soundboard-api/` — Go API + SQLite (`modernc.org/sqlite`, no cgo) backing the soundboard page, plus a Go CLI for uploading new clips. The frontend fetches from this API client-side rather than at build time, since the clip list changes without a redeploy.

Deployment target is tentatively Cloudflare Pages (`adapter-cloudflare`) for the frontend — not finalized. Backend hosting, and where audio files themselves live (on-disk vs. Cloudflare R2), are also still open.

## Pages

- Home / Links
- About
- Ref Sheet (image gallery)
- Credits
- Soundboard (grid of clickable clips, backed by the Go API)

## Status

In progress.

- `web/`: SvelteKit skeleton scaffolded and running (`npm run dev`). No real pages/routes built yet.
- `soundboard-api/`: Go module scaffolded (`cmd/api`, `cmd/cli`, `internal/soundbite`, `internal/db`, `internal/httpapi`, `internal/storage`). No real logic yet — SQLite schema and HTTP server are next.

## Planned / undecided

- Static content (links, about text, credits, ref sheet entries) will live in Markdown/JSON files, not hardcoded in components.
- Soundbite schema still needs tags/category and source/stream context, deferred until there are more real clips to design against.
- Final deployment targets for the frontend and backend, and audio file storage location.
