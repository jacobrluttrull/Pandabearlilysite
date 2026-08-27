import type { Soundbite } from './types';

export async function fetchSoundbites(
	baseUrl: string,
	fetchFn: typeof fetch = fetch
): Promise<Soundbite[]> {
	const response = await fetchFn(`${baseUrl}/soundbites`);

	if (!response.ok) {
		throw new Error(`failed to fetch soundbites: ${response.status} ${response.statusText}`);
	}

	return response.json();
}

/**
 * Records one play of a clip and resolves to the clip's new total.
 *
 * Sent as a bare POST with no body or custom headers, which keeps it a "simple" CORS
 * request — no preflight round trip ahead of every tap.
 */
export async function recordPlay(
	baseUrl: string,
	id: number,
	fetchFn: typeof fetch = fetch
): Promise<number> {
	const response = await fetchFn(`${baseUrl}/soundbites/${id}/play`, { method: 'POST' });

	if (!response.ok) {
		throw new Error(`failed to record play: ${response.status} ${response.statusText}`);
	}

	const body = await response.json();
	return body.play_count;
}
