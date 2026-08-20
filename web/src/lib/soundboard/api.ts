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
