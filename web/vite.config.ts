import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},

			// Every page is prerendered to plain HTML, which the Go API serves alongside
			// the soundboard endpoints. One process, one origin, no CORS.
			adapter: adapter({
				pages: 'build',
				assets: 'build',
				// Prerendering covers every route, so the SPA fallback is only ever hit by
				// an unknown URL — the Go handler turns that into a real 404.
				fallback: '404.html',
				precompress: false,
				strict: true
			})
		})
	],
	server: {
		// In dev the site and API are separate processes on different ports. Proxying
		// /api keeps the frontend's fetch paths relative, exactly as they are in
		// production, so there is no environment-specific URL to get wrong.
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	test: {
		environment: 'node'
	}
});
