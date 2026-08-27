<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchSoundbites, recordPlay } from '$lib/soundboard/api';
	import type { Soundbite } from '$lib/soundboard/types';
	import Button from '$lib/components/Button.svelte';
	import PawPrint from '$lib/components/PawPrint.svelte';
	import BambooStalk from '$lib/components/BambooStalk.svelte';

	// Same-origin in production; in dev, Vite proxies /api to the Go process. Keeping
	// the path relative means there is no environment-specific URL to configure.
	const baseUrl = '';

	type Status = 'loading' | 'error' | 'ready';

	let status = $state<Status>('loading');
	let soundbites = $state<Soundbite[]>([]);
	let query = $state('');
	let nowPlayingId = $state<number | null>(null);
	let audioEl = $state<HTMLAudioElement>();

	// A full board is a wall of tiles on arrival, so it opens on a partial grid.
	// Searching narrows the list on its own, so the cap steps out of the way.
	const INITIAL_VISIBLE = 24;
	let showAll = $state(false);

	const filtered = $derived.by(() => {
		const q = query.trim().toLowerCase();
		return q ? soundbites.filter((clip) => clip.name.toLowerCase().includes(q)) : soundbites;
	});

	const searching = $derived(query.trim().length > 0);
	const collapsed = $derived(!showAll && !searching && filtered.length > INITIAL_VISIBLE);
	const visible = $derived(collapsed ? filtered.slice(0, INITIAL_VISIBLE) : filtered);
	const restCount = $derived(filtered.length - visible.length);

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
		countPlay(clip);
	}

	// Clips run a second or two, so pressing the tile is the play — the tally is not
	// held back waiting for playback to finish.
	async function countPlay(clip: Soundbite) {
		// Move the number immediately so the tile feels responsive, then reconcile with
		// the server's real total once the request lands.
		clip.play_count += 1;
		try {
			clip.play_count = await recordPlay(baseUrl, clip.id);
		} catch {
			// A lost tally is not worth interrupting playback over; the optimistic
			// bump stands until the next page load corrects it.
		}
	}

	function handleEnded() {
		nowPlayingId = null;
	}

	// Keeps long tallies from widening a tile: 999 -> "999", 1240 -> "1.2k".
	function formatCount(plays: number): string {
		if (plays < 1000) return String(plays);
		if (plays < 10_000) return `${(plays / 1000).toFixed(1)}k`;
		if (plays < 1_000_000) return `${Math.round(plays / 1000)}k`;
		return `${(plays / 1_000_000).toFixed(1)}m`;
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

	<div class="search-row">
		<label class="search">
			<svg class="search-icon" viewBox="0 0 24 24" width="16" height="16" aria-hidden="true">
				<circle cx="10.5" cy="10.5" r="6.5" fill="none" stroke="currentColor" stroke-width="2" />
				<line x1="15.3" y1="15.3" x2="21" y2="21" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
			</svg>
			<span class="sr-only">Search sounds</span>
			<input type="search" placeholder="Search sounds…" bind:value={query} />
		</label>

		{#if status === 'ready' && soundbites.length > 0}
			<p class="clip-count" aria-live="polite">
				{searching
					? `${filtered.length} of ${soundbites.length} clips`
					: `${soundbites.length} clips`}
			</p>
		{/if}
	</div>

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
		<div class="grid" class:collapsed>
			{#each visible as clip (clip.id)}
				{@const playing = nowPlayingId === clip.id}
				<!-- the tile is a button, so the download link sits beside it in a
				     wrapper rather than inside it — nesting an <a> in a <button> is
				     invalid and breaks keyboard navigation -->
				<div class="clip-cell">
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
						<span class="meta-stats">
							<span class="plays">
								<svg viewBox="0 0 12 12" width="9" height="9" aria-hidden="true">
									<path d="M2.5 1.3 10 6l-7.5 4.7z" fill="currentColor" />
								</svg>
								{formatCount(clip.play_count)}
								<span class="sr-only">
									{clip.play_count === 1 ? 'play' : 'plays'}
								</span>
							</span>
							<span class="duration">{formatDuration(clip.length_seconds)}</span>
						</span>
					</span>
				</button>
					<a
						class="clip-download"
						href={`${baseUrl}${clip.download_url}`}
						download
						title={`Download "${clip.name}"`}
					>
						<svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
							<path
								d="M8 1.8v7.2m0 0L5.1 6.1M8 9L10.9 6.1M2.7 11.4v1.6a1.2 1.2 0 0 0 1.2 1.2h8.2a1.2 1.2 0 0 0 1.2-1.2v-1.6"
								fill="none"
								stroke="currentColor"
								stroke-width="1.6"
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
						</svg>
						<span class="sr-only">Download {clip.name}</span>
					</a>
				</div>
			{/each}
		</div>

		{#if !searching && filtered.length > INITIAL_VISIBLE}
			<div class="more-row">
				<Button
					variant="secondary"
					size="sm"
					aria-expanded={showAll}
					onclick={() => (showAll = !showAll)}
				>
					{showAll ? 'Show fewer clips' : `Show all ${filtered.length} clips`}
				</Button>
				{#if collapsed}
					<span class="more-hint">{restCount} more waiting in the grove</span>
				{/if}
			</div>
		{/if}
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
	.search-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-2);
		margin-bottom: var(--space-3);
	}

	.search {
		position: relative;
		display: block;
		flex: 1 1 18rem;
		max-width: 26rem;
	}

	/* orientation: how big the board is, and how much of it the search left */
	.clip-count {
		margin: 0;
		color: var(--color-text-muted);
		font: var(--text-label-sm);
		letter-spacing: 0.02em;
		font-variant-numeric: tabular-nums;
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
		grid-template-columns: repeat(auto-fill, minmax(var(--tile-min-w), 1fr));
		gap: var(--space-1);
	}

	/* holds the play tile and its download link as siblings, so the link can sit
	   over the tile's corner without nesting interactive elements */
	.clip-cell {
		position: relative;
		display: flex;
	}

	/* A tile rests quiet — soft fill, hairline edge — and only takes on the
	   theme's full surface treatment when it's hovered, focused, or sounding.
	   56 tiles at full weight is the wall; 56 tiles at rest is a field. */
	.clip-tile {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-1);
		min-height: var(--tile-min-h);
		padding: var(--space-2);
		border-radius: var(--radius);
		border: 1px solid var(--surface-quiet-border);
		background: var(--surface-quiet-bg);
		backdrop-filter: var(--surface-blur);
		-webkit-backdrop-filter: var(--surface-blur);
		color: var(--color-text);
		text-align: left;
		cursor: pointer;
		font: var(--text-body-md);
		transition:
			transform 0.15s ease,
			border-color 0.2s ease,
			background-color 0.2s ease,
			box-shadow 0.2s ease;
	}

	/* the name carries the tile — it wraps to as many lines as it needs
	   rather than truncating mid-word. The right inset keeps it clear of the
	   download link sitting in the corner. */
	.clip-name {
		font-weight: 600;
		line-height: 1.3;
		overflow-wrap: break-word;
		padding-right: var(--space-3);
	}

	/* Always visible rather than hover-only: on a touch screen there is no hover,
	   and a save action nobody can find is not a feature. It stays muted so it
	   reads as secondary to the tile's main job, which is playing. */
	.clip-download {
		position: absolute;
		top: var(--space-1);
		right: var(--space-1);
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1.55rem;
		height: 1.55rem;
		border-radius: var(--radius-sm);
		color: var(--color-text-muted);
		opacity: 0.75;
		transition:
			color 0.2s ease,
			background-color 0.2s ease,
			opacity 0.2s ease;
	}

	.clip-download:hover,
	.clip-download:focus-visible {
		color: var(--color-accent);
		background: var(--surface-quiet-bg);
		opacity: 1;
	}

	.clip-download:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 1px;
	}

	/* while a clip sounds, its download follows the tile into canopy gold */
	.clip-cell:has(.clip-tile.playing) .clip-download {
		color: var(--color-canopy-strong);
		opacity: 1;
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

	.meta-stats {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		min-width: 0;
	}

	.duration {
		font: var(--text-label-sm);
		color: var(--color-text-muted);
		font-variant-numeric: tabular-nums;
	}

	.plays {
		display: inline-flex;
		align-items: center;
		gap: 0.3em;
		font: var(--text-label-sm);
		color: var(--color-text-muted);
		font-variant-numeric: tabular-nums;
	}

	/* the tally leans into the canopy gold while its clip is sounding, so the number
	   reads as "this just went up" rather than as static metadata */
	.clip-tile.playing .plays {
		color: var(--color-canopy-strong);
	}

	/* the waveform only speaks while a clip is sounding — 56 dormant glyphs is
	   noise, so an idle tile holds the space without spending the ink */
	.clip-tile .bars {
		opacity: 0;
		transition: opacity 0.2s ease;
	}

	.clip-tile.playing .bars {
		opacity: 1;
	}

	.clip-tile:hover:not(.playing),
	.clip-tile:focus-visible:not(.playing) {
		background: var(--surface-glass-bg);
		border-color: var(--color-accent);
		transform: translateY(-2px);
		box-shadow: 0 0 0.9rem color-mix(in srgb, var(--color-accent) 30%, transparent);
	}

	/* metadata sits back at rest and steps up when the tile is engaged */
	.clip-tile:hover:not(.playing) :is(.plays, .duration),
	.clip-tile:focus-visible:not(.playing) :is(.plays, .duration) {
		color: var(--color-text);
	}

	.clip-tile:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	/* revealing the rest of the board */
	.more-row {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--space-2);
		margin-top: var(--space-3);
	}

	.more-hint {
		color: var(--color-text-muted);
		font: var(--text-label-sm);
		font-variant-numeric: tabular-nums;
	}

	/* signature move: a tile that's actually playing comes alive — the paw
	   glyph steps in, the waveform snaps into an equalizer, and the whole
	   tile breathes between canopy gold and rust glow */
	.clip-tile.playing {
		background: var(--surface-glass-bg);
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

	/* reduced motion keeps every state cue — gold border, lit waveform, paw —
	   and drops only the looping movement */
	@media (prefers-reduced-motion: reduce) {
		.clip-tile,
		.clip-tile.playing,
		.clip-tile:active,
		.bars.animate i,
		.clip-name :global(.playing-paw) {
			animation: none;
			transition: none;
		}

		.clip-tile.playing {
			box-shadow: 0 0 0 1px var(--color-canopy-strong) inset;
		}

		.bars.animate i {
			height: 100%;
		}
	}

	@media (max-width: 480px) {
		.header-row {
			align-items: flex-start;
		}

		:global(.header-stalk) {
			display: none;
		}

		/* a phone still holds two columns — one column of tall tiles turns the
		   board into a very long scroll */
		.grid {
			grid-template-columns: repeat(auto-fill, minmax(var(--tile-min-w-sm), 1fr));
		}
	}
</style>
