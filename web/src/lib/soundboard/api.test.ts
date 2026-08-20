import { describe, expect, it, vi } from 'vitest';
import { fetchSoundbites } from './api';
import type { Soundbite } from './types';

const sampleSoundbites: Soundbite[] = [
	{
		id: 1,
		name: 'Good Morning',
		date_made: '2026-01-15',
		date_stored: '2026-08-20T00:14:50.862Z',
		length_seconds: 1.8,
		audio_url: '/soundbites/1/audio'
	}
];

function fakeFetch(body: unknown, ok = true, status = 200): typeof fetch {
	return vi.fn().mockResolvedValue({
		ok,
		status,
		statusText: ok ? 'OK' : 'Internal Server Error',
		json: () => Promise.resolve(body)
	}) as unknown as typeof fetch;
}

describe('fetchSoundbites', () => {
	it('returns the parsed list of soundbites', async () => {
		const fetchFn = fakeFetch(sampleSoundbites);

		const result = await fetchSoundbites('http://localhost:8080', fetchFn);

		expect(result).toEqual(sampleSoundbites);
	});

	it('requests the /soundbites endpoint on the given base URL', async () => {
		const fetchFn = fakeFetch(sampleSoundbites);

		await fetchSoundbites('http://localhost:8080', fetchFn);

		expect(fetchFn).toHaveBeenCalledWith('http://localhost:8080/soundbites');
	});

	it('returns an empty array when there are no soundbites yet', async () => {
		const fetchFn = fakeFetch([]);

		const result = await fetchSoundbites('http://localhost:8080', fetchFn);

		expect(result).toEqual([]);
	});

	it('throws when the API responds with an error status', async () => {
		const fetchFn = fakeFetch({ error: 'boom' }, false, 500);

		await expect(fetchSoundbites('http://localhost:8080', fetchFn)).rejects.toThrow(/500/);
	});
});
