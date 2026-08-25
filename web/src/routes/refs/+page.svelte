<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import refSheet from '$lib/data/ref-sheet.json';

	type RefEntry = {
		id: string;
		title: string;
		caption: string;
		image: string;
	};

	const entries = refSheet as RefEntry[];

	// tracks which entries' <img> failed to load, so we can fall back to
	// the placeholder box instead of showing a broken-image icon
	let failed = $state<Record<string, boolean>>({});

	function markFailed(id: string) {
		failed[id] = true;
	}

	function showsPlaceholder(entry: RefEntry) {
		return !entry.image || failed[entry.id];
	}

	let dialogEl: HTMLDialogElement;
	let active = $state<RefEntry | null>(null);

	function openEntry(entry: RefEntry) {
		active = entry;
		dialogEl.showModal();
	}

	function closeDialog() {
		dialogEl.close();
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === dialogEl) {
			closeDialog();
		}
	}
</script>

<svelte:head>
	<title>PandaLily — Ref Sheet</title>
</svelte:head>

<section class="refs-page">
	<header class="refs-header">
		<h1>Reference Sheet</h1>
		<p>Official visual reference for PandaLily — body, expressions, colors, and outfits.</p>
	</header>

	<div class="gallery">
		{#each entries as entry (entry.id)}
			<Card>
				<button type="button" class="gallery-trigger" onclick={() => openEntry(entry)}>
					<span class="frame">
						{#if showsPlaceholder(entry)}
							<span class="placeholder">
								<svg
									viewBox="0 0 24 24"
									width="32"
									height="32"
									fill="none"
									stroke="currentColor"
									stroke-width="1.5"
									aria-hidden="true"
								>
									<rect x="3" y="4" width="18" height="16" rx="2" />
									<circle cx="8.5" cy="10" r="1.75" />
									<path d="M21 16.5 15.5 11 5 20" />
								</svg>
								<span class="placeholder-label">Image coming soon</span>
							</span>
						{:else}
							<img src={entry.image} alt={entry.title} onerror={() => markFailed(entry.id)} />
						{/if}
					</span>
					<span class="entry-body">
						<span class="entry-title">{entry.title}</span>
						<span class="entry-caption">{entry.caption}</span>
					</span>
				</button>
			</Card>
		{/each}
	</div>
</section>

<dialog bind:this={dialogEl} class="lightbox" onclick={handleBackdropClick} onclose={() => (active = null)}>
	{#if active}
		{@const entry = active}
		<div class="lightbox-content">
			<button type="button" class="lightbox-close" onclick={closeDialog} aria-label="Close">✕</button>
			<span class="frame lightbox-frame">
				{#if showsPlaceholder(entry)}
					<span class="placeholder">
						<svg
							viewBox="0 0 24 24"
							width="48"
							height="48"
							fill="none"
							stroke="currentColor"
							stroke-width="1.5"
							aria-hidden="true"
						>
							<rect x="3" y="4" width="18" height="16" rx="2" />
							<circle cx="8.5" cy="10" r="1.75" />
							<path d="M21 16.5 15.5 11 5 20" />
						</svg>
						<span class="placeholder-label">Image coming soon</span>
					</span>
				{:else}
					<img src={entry.image} alt={entry.title} onerror={() => markFailed(entry.id)} />
				{/if}
			</span>
			<div class="lightbox-body">
				<h2>{entry.title}</h2>
				<p>{entry.caption}</p>
			</div>
		</div>
	{/if}
</dialog>

<style>
	.refs-page {
		max-width: var(--container-max);
		margin: 0 auto;
		padding: var(--space-4) var(--space-2) var(--space-6);
	}

	.refs-header {
		max-width: 640px;
		margin: 0 auto var(--space-4);
		text-align: center;
	}

	.refs-header h1 {
		font: var(--text-headline-lg-mobile);
		margin-bottom: var(--space-1);
	}

	.refs-header p {
		margin: 0;
		color: var(--color-text-muted);
		font: var(--text-body-lg);
	}

	.gallery {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-3);
	}

	.gallery-trigger {
		display: flex;
		flex-direction: column;
		width: 100%;
		border: none;
		background: none;
		padding: 0;
		margin: 0;
		font: inherit;
		color: inherit;
		text-align: left;
		cursor: pointer;
		border-radius: inherit;
	}

	.gallery-trigger:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: -2px;
	}

	.frame {
		display: block;
		width: 100%;
		aspect-ratio: 4 / 5;
		background: var(--color-surface-alt);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		overflow: hidden;
	}

	.frame img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.3s ease;
	}

	.gallery-trigger:hover .frame img {
		transform: scale(1.03);
	}

	.placeholder {
		display: flex;
		height: 100%;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-1);
		color: var(--color-text-muted);
	}

	.placeholder-label {
		font: var(--text-label-sm);
		text-transform: uppercase;
		letter-spacing: 0.06em;
	}

	.entry-body {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding-top: var(--space-2);
	}

	.entry-title {
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.05rem;
		color: var(--color-text);
	}

	.entry-caption {
		font: var(--text-body-md);
		color: var(--color-text-muted);
	}

	@media (min-width: 640px) {
		.gallery {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (min-width: 1024px) {
		.gallery {
			grid-template-columns: repeat(3, 1fr);
		}
	}

	/* lightbox */
	.lightbox {
		max-width: min(90vw, 720px);
		width: 100%;
		border: none;
		border-radius: var(--radius-lg);
		padding: 0;
		background: var(--color-surface);
		color: var(--color-text);
	}

	.lightbox::backdrop {
		background: color-mix(in srgb, var(--color-ink) 60%, transparent);
		backdrop-filter: blur(4px);
	}

	.lightbox-content {
		position: relative;
		display: flex;
		flex-direction: column;
	}

	.lightbox-close {
		position: absolute;
		top: var(--space-2);
		right: var(--space-2);
		z-index: 1;
		width: 2rem;
		height: 2rem;
		border-radius: var(--radius-full);
		border: none;
		background: color-mix(in srgb, var(--color-ink) 55%, transparent);
		color: var(--color-on-ink);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.lightbox-frame {
		aspect-ratio: 4 / 5;
		border-radius: var(--radius-lg) var(--radius-lg) 0 0;
	}

	.lightbox-body {
		padding: var(--space-3);
	}

	.lightbox-body h2 {
		margin-bottom: 0.35rem;
	}

	.lightbox-body p {
		margin: 0;
		color: var(--color-text-muted);
	}
</style>
