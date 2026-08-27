// The clip list is fetched in the browser, not at build time, so the soundboard can
// prerender to a static shell like every other page. What ships is the loading state;
// onMount fills in the grid from the API.
export const prerender = true;
