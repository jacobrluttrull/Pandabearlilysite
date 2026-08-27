// Every page is prerendered to plain HTML at build time. The whole site is static
// content — the only live data is the soundboard's clip list, which is fetched in the
// browser rather than at build time, so it prerenders to a shell like everything else.
//
// This is what lets the Go API serve the site directly: no Node server, no SSR, just
// files on disk next to the audio.
export const prerender = true;
