<script lang="ts">
	import { onMount } from 'svelte';
	import { env } from '$env/dynamic/public';
	import { fetchSoundbites } from '$lib/soundboard/api';
	import type { Soundbite } from '$lib/soundboard/types';
	import Button from '$lib/components/Button.svelte';

	const baseUrl = env.PUBLIC_SOUNDBOARD_API_URL || 'http://localhost:8080';

	type Status = 'loading' | 'error' | 'ready';

	let status = $state<Status>('loading');
	let soundbites = $state<Soundbite[]>([]);
	let query = $state('');
	let nowPlayingId = $state<number | null>(null);
	let audioEl = $state<HTMLAudioElement>();

	const filtered = $derived.by(() => {
		const q = query.trim().toLowerCase();
		return q ? soundbites.filter((clip) => clip.name.toLowerCase().includes(q)) : soundbites;
	});

	onMount(async () => {
		try {
			soundbites = await fetchSoundbites(baseUrl);
			status = 'ready';
		} catch {
			status = 'error';
		}
	});

	function play(clip: Soundbite) {
		if (!audioEl) return;
		audioEl.src = `${baseUrl}${clip.audio_url}`;
		nowPlayingId = clip.id;
		audioEl.play();
	}

	function handleEnded() {
		nowPlayingId = null;
	}
</script>

<svelte:head>
	<title>PandaLily — Soundboard</title>
</svelte:head>

<section class="soundboard">
	<h1>Soundboard</h1>

	<div class="search">
		<input
			type="search"
			placeholder="Search sounds…"
			aria-label="Search sounds"
			bind:value={query}
		/>
	</div>

	{#if status === 'loading'}
		<p class="status-message">Loading clips…</p>
	{:else if status === 'error'}
		<p class="status-message">
			Couldn't reach the soundboard right now. If you're running this locally, make sure the API
			is up.
		</p>
	{:else if filtered.length === 0}
		<p class="status-message">
			{soundbites.length === 0 ? 'No clips yet — check back soon.' : 'No clips match your search.'}
		</p>
	{:else}
		<div class="grid">
			{#each filtered as clip (clip.id)}
				<Button
					type="button"
					variant={nowPlayingId === clip.id ? 'primary' : 'secondary'}
					class="clip"
					onclick={() => play(clip)}
				>
					{clip.name}
				</Button>
			{/each}
		</div>
	{/if}
</section>

<audio bind:this={audioEl} onended={handleEnded}></audio>

<style>
	.soundboard {
		max-width: var(--container-max);
		margin: 0 auto;
		padding: var(--space-4) var(--space-3) var(--space-6);
	}

	h1 {
		font: var(--text-headline-lg);
		margin: 0 0 var(--space-3);
	}

	.search {
		margin-bottom: var(--space-4);
	}

	.search input {
		width: 100%;
		max-width: 32rem;
		padding: var(--space-2) var(--space-3);
		border-radius: var(--radius-full);
		border: 1px solid var(--color-border);
		background: var(--color-surface);
		color: var(--color-text);
		font: var(--text-body-md);
	}

	.search input:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	.status-message {
		color: var(--color-text-muted);
		font: var(--text-body-lg);
		padding: var(--space-4) 0;
	}

	.grid {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-2);
	}

	:global(.clip) {
		text-align: center;
	}
</style>
