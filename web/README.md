# web

The PandaLily fan site: SvelteKit 2, Svelte 5 in runes mode, TypeScript, and plain CSS.
Six pages — Home/Links, About, Art, Ref Sheet, Credits, and the Soundboard.

Every route is prerendered to plain HTML by `adapter-static`, so what ships is a folder of
files. In production the Go binary in [`../soundboard-api`](../soundboard-api) serves those
files and the API from one process; there is no Node server at runtime.

## Running it

```powershell
npm install
npm run dev          # http://localhost:5173
```

The soundboard needs the API running alongside it — `go run .\cmd\api` from
`../soundboard-api`, listening on `:8080`. Vite proxies `/api` there (see
`vite.config.ts`), which is what lets the frontend use the same relative paths in
development that it uses in production. **There are no environment variables to set**;
`.env.example` exists only as a placeholder in case a real public one is ever needed.

| Command | What it does |
| --- | --- |
| `npm run dev` | Dev server with HMR |
| `npm run build` | Prerender the whole site into `build/` |
| `npm run preview` | Serve the built output |
| `npm run check` | `svelte-check` against `tsconfig.json` |
| `npm test` | Vitest run |

## Layout

```
src/
  routes/            → one directory per page; +layout.ts sets prerender = true
  lib/
    data/            → page content as JSON (links, about, art, credits, ref-sheet)
    components/      → Button, Card, Chip, Nav, Footer, Gallery, icons, decorations
    soundboard/      → API client and types for the clip grid
    styles/          → tokens.css (theme variables) and global.css
static/              → images, favicon, robots.txt
```

**Content lives in `src/lib/data/*.json`, not in components.** Adding a link, an artwork, a
credit, or a ref sheet entry is a JSON edit; the page reads the file and renders it.

**Styling goes through tokens.** Every color, space, radius, and font is a CSS custom
property in `src/lib/styles/tokens.css`, defined per theme and selected by `data-theme` on
`<html>`. Components reference `var(--token)` and never hardcode a value, so a theme tweak
happens in one file. The two themes are `bamboo` (dark, glassmorphic — the default) and
`cream` (light, papery); the nav's toggle switches them and stores the choice in
`localStorage`, and an inline script in `app.html` applies the stored choice before first
paint so the page never flashes the wrong theme. The full pattern, including the checklist
for new components, is in [`../CLAUDE.md`](../CLAUDE.md).

## The soundboard page

It is the one page whose content is not known at build time, because clips are added
without redeploying the site. So it prerenders to a loading shell and fetches
`/api/soundbites` in `onMount`, then renders the grid: tap a tile to play, a search box to
filter, a play tally per clip that bumps optimistically and reconciles with the server, and
a download button per clip. The board opens capped at 24 tiles so a full collection is not
a wall on arrival; searching lifts the cap.
