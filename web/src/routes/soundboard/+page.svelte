<script lang="ts">
	import { onMount } from 'svelte';
	import { env } from '$env/dynamic/public';
	import { fetchSoundbites } from '$lib/soundboard/api';
	import type { Soundbite } from '$lib/soundboard/types';
	import PawPrint from '$lib/components/PawPrint.svelte';
	import BambooStalk from '$lib/components/BambooStalk.svelte';

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

	function formatDuration(totalSeconds: number): string {
		const s = Math.max(0, Math.round(totalSeconds));
		const m = Math.floor(s / 60);
		const rem = s % 60;
		return `${m}:${rem.toString().padStart(2, '0')}`;
	}
</script>

<svelte:head>
	<title>PandaLily — Soundboard</title>
</svelte:head>

<section class="soundboard">
	<div class="header-row">
		<div>
			<h1>Soundboard</h1>
			<p class="subtitle">Tap a clip to hear her voice ring through the grove.</p>
		</div>
		<BambooStalk segments={2} leaf sway class="header-stalk" aria-hidden="true" />
	</div>

	<label class="search">
		<svg class="search-icon" viewBox="0 0 24 24" width="16" height="16" aria-hidden="true">
			<circle cx="10.5" cy="10.5" r="6.5" fill="none" stroke="currentColor" stroke-width="2" />
			<line x1="15.3" y1="15.3" x2="21" y2="21" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
		</svg>
		<span class="sr-only">Search sounds</span>
		<input type="search" placeholder="Search sounds…" bind:value={query} />
	</label>

	{#if status === 'loading'}
		<p class="status-message">
			<span class="bars animate loading-bars" aria-hidden="true"><i></i><i></i><i></i></span>
			Tuning the chimes…
		</p>
	{:else if status === 'error'}
		<p class="status-message">The soundboard is quiet right now — check back in a bit.</p>
	{:else if filtered.length === 0}
		<p class="status-message">
			{soundbites.length === 0
				? 'No chimes strung up yet — check back soon.'
				: `Nothing rings back for "${query}" — try another search.`}
		</p>
	{:else}
		<div class="grid">
			{#each filtered as clip (clip.id)}
				{@const playing = nowPlayingId === clip.id}
				<button
					type="button"
					class="clip-tile"
					class:playing
					title={clip.name}
					aria-pressed={playing}
					onclick={() => play(clip)}
				>
					<span class="clip-name">
						{#if playing}
							<PawPrint size={13} class="playing-paw" />
						{/if}
						{clip.name}
					</span>
					<span class="clip-meta">
						<span class="bars" class:animate={playing} aria-hidden="true">
							<i></i><i></i><i></i>
						</span>
						<span class="duration">{formatDuration(clip.length_seconds)}</span>
					</span>
				</button>
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

	.header-row {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--space-3);
		margin-bottom: var(--space-4);
	}

	h1 {
		font: var(--text-headline-lg);
		margin: 0 0 var(--space-1);
	}

	.subtitle {
		margin: 0;
		color: var(--color-text-muted);
		font: var(--text-body-lg);
		max-width: 38ch;
	}

	:global(.header-stalk) {
		flex-shrink: 0;
		color: var(--color-accent);
		margin-bottom: var(--space-1);
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}

	/* search */
	.search {
		position: relative;
		display: block;
		width: 100%;
		max-width: 26rem;
		margin-bottom: var(--space-4);
	}

	.search-icon {
		position: absolute;
		top: 50%;
		left: var(--space-2);
		transform: translateY(-50%);
		color: var(--color-text-muted);
		pointer-events: none;
	}

	.search input {
		width: 100%;
		padding: var(--space-2) var(--space-3) var(--space-2) calc(var(--space-2) * 2 + 1rem);
		border-radius: var(--radius-full);
		border: 1px solid var(--color-border);
		background: var(--color-surface);
		color: var(--color-text);
		font: var(--text-body-md);
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease;
	}

	.search input::placeholder {
		color: var(--color-text-muted);
	}

	.search input:focus-visible {
		outline: none;
		border-color: var(--color-canopy-strong);
		box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-canopy) 30%, transparent);
	}

	/* status states */
	.status-message {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		color: var(--color-text-muted);
		font: var(--text-body-lg);
		padding: var(--space-4) 0;
	}

	/* waveform glyph — dormant on idle tiles, wakes into an equalizer while
	   a clip plays (and doubles as the loading state's "tuning up" cue) */
	.bars {
		display: inline-flex;
		align-items: flex-end;
		gap: 2px;
		height: 0.85rem;
		flex-shrink: 0;
	}

	.bars i {
		display: block;
		width: 3px;
		height: 30%;
		border-radius: 1px;
		background: var(--color-text-muted);
		opacity: 0.55;
		transition:
			height 0.2s ease,
			opacity 0.2s ease,
			background-color 0.2s ease;
	}

	.bars.animate i {
		opacity: 1;
		animation: bar-bounce 0.9s ease-in-out infinite;
	}

	.bars.animate i:nth-child(1) {
		animation-delay: 0s;
	}
	.bars.animate i:nth-child(2) {
		animation-delay: 0.15s;
	}
	.bars.animate i:nth-child(3) {
		animation-delay: 0.3s;
	}

	.loading-bars i {
		background: var(--color-accent);
	}

	@keyframes bar-bounce {
		0%,
		100% {
			height: 30%;
		}
		50% {
			height: 100%;
		}
	}

	/* clip grid */
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(9.5rem, 1fr));
		gap: var(--space-2);
	}

	.clip-tile {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-2);
		min-height: 5.5rem;
		padding: var(--space-2) var(--space-3);
		border-radius: var(--radius-lg);
		border: 1px solid var(--surface-glass-border);
		background: var(--surface-glass-bg);
		backdrop-filter: var(--surface-blur);
		-webkit-backdrop-filter: var(--surface-blur);
		color: var(--color-text);
		text-align: left;
		cursor: pointer;
		font: var(--text-body-md);
		transition:
			transform 0.15s ease,
			border-color 0.2s ease,
			box-shadow 0.2s ease;
	}

	.clip-name {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
		font-weight: 600;
		line-height: 1.3;
	}

	.clip-name :global(.playing-paw) {
		display: inline-block;
		color: var(--color-panda-rust);
		margin-right: 0.3em;
		vertical-align: -0.05em;
		animation: paw-pop 0.3s ease-out;
	}

	.clip-meta {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-2);
		width: 100%;
	}

	.duration {
		font: var(--text-label-sm);
		color: var(--color-text-muted);
		font-variant-numeric: tabular-nums;
	}

	.clip-tile:hover:not(.playing) {
		border-color: var(--color-accent);
		transform: translateY(-2px);
		box-shadow: 0 0 0.9rem color-mix(in srgb, var(--color-accent) 30%, transparent);
	}

	.clip-tile:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	/* signature move: a tile that's actually playing comes alive — the paw
	   glyph steps in, the waveform snaps into an equalizer, and the whole
	   tile breathes between canopy gold and rust glow */
	.clip-tile.playing {
		border-color: var(--color-canopy-strong);
		animation: chime-glow 1.6s ease-in-out infinite;
	}

	.clip-tile.playing .bars.animate i {
		background: var(--color-canopy-strong);
	}

	@keyframes chime-glow {
		0%,
		100% {
			transform: translateY(-1px);
			box-shadow:
				0 0 0.75rem color-mix(in srgb, var(--color-canopy) 45%, transparent),
				0 0 0 1px var(--color-canopy-strong) inset;
		}
		50% {
			transform: translateY(-2px) scale(1.015);
			box-shadow:
				0 0 1.4rem color-mix(in srgb, var(--color-panda-rust) 50%, transparent),
				0 0 0 1px var(--color-canopy-strong) inset;
		}
	}

	@keyframes paw-pop {
		0% {
			transform: scale(0.4);
			opacity: 0;
		}
		60% {
			transform: scale(1.15);
			opacity: 1;
		}
		100% {
			transform: scale(1);
		}
	}

	/* tactile press feedback — a quick outward burst on every tap,
	   independent of the sustained "now playing" glow above */
	.clip-tile:active {
		animation: chime-press 0.35s ease-out;
	}

	@keyframes chime-press {
		0% {
			transform: scale(0.96);
			box-shadow: 0 0 0 0 color-mix(in srgb, var(--color-accent) 55%, transparent);
		}
		100% {
			transform: scale(1);
			box-shadow: 0 0 0 0.85rem color-mix(in srgb, var(--color-accent) 0%, transparent);
		}
	}

	@media (max-width: 480px) {
		.header-row {
			align-items: flex-start;
		}

		:global(.header-stalk) {
			display: none;
		}
	}
</style>
