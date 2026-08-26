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
- Art (commissioned artwork gallery)
- Ref Sheet (image gallery)
- Credits
- Soundboard (grid of clickable clips, backed by the Go API)

## Design system

Two themes — "bamboo" (dark, glassmorphic, default) and "cream" (light, flat) — selected via `data-theme` on `<html>` and stored as a user preference. Every color, spacing, radius, and font value is a CSS custom property in `web/src/lib/styles/tokens.css`; components (`Button`, `Card`, `Chip`, `Nav`, `Gallery`, etc.) reference tokens only, never hardcoded values. See the "PandaLily Design Pattern" section in `CLAUDE.md` for the full pattern.

## Status

In progress.

- `web/`: All six pages built and styled with the design system above. Static content (links, about text, art entries, credits, ref sheet entries) lives in `web/src/lib/data/*.json`, not hardcoded in components.
- `soundboard-api/`: Go API (`cmd/api`) serving the soundbite list and audio files, SQLite (`modernc.org/sqlite`) with sqlc-generated queries and migrations, and a CLI (`cmd/cli`) for uploading and listing clips. The frontend soundboard page fetches from this API client-side (not at build time).

## Planned / undecided

- Soundbite schema still needs tags/category and source/stream context, deferred until there are more real clips to design against.
- Final deployment targets for the frontend and backend, and audio file storage location (on-disk next to the API vs. Cloudflare R2).
