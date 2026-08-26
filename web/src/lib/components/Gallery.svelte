<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import PawPrint from '$lib/components/PawPrint.svelte';
	import BambooStalk from '$lib/components/BambooStalk.svelte';
	import BambooSprout from '$lib/components/BambooSprout.svelte';

	type GalleryEntry = {
		id: string;
		title: string;
		caption: string;
		image: string;
		url?: string;
	};

	type GalleryProps = {
		heading: string;
		subtitle: string;
		entries: GalleryEntry[];
		linkLabel?: string;
	};

	let { heading, subtitle, entries, linkLabel = 'Visit portfolio' }: GalleryProps = $props();

	// tracks which entries' <img> failed to load, so we can fall back to
	// the placeholder box instead of showing a broken-image icon
	let failed = $state<Record<string, boolean>>({});

	function markFailed(id: string) {
		failed[id] = true;
	}

	function showsPlaceholder(entry: GalleryEntry) {
		return !entry.image || failed[entry.id];
	}

	let dialogEl: HTMLDialogElement;
	let active = $state<GalleryEntry | null>(null);

	function openEntry(entry: GalleryEntry) {
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

<section class="gallery-page">
	<header class="gallery-header">
		<h1>{heading}</h1>
		<p>{subtitle}</p>
		<div class="grove-divider" aria-hidden="true">
			<BambooStalk segments={2} leaf={false} class="grove-divider-stalk" />
		</div>
	</header>

	<div class="gallery">
		{#each entries as entry (entry.id)}
			<Card>
				<button type="button" class="gallery-trigger" onclick={() => openEntry(entry)}>
					<span class="frame">
						{#if showsPlaceholder(entry)}
							<span class="mat placeholder">
								<BambooSprout size={40} class="placeholder-sprout" />
								<span class="placeholder-label">Image coming soon</span>
							</span>
						{:else}
							<span class="mat">
								<img src={entry.image} alt={entry.title} onerror={() => markFailed(entry.id)} />
							</span>
						{/if}
					</span>
					<span class="entry-body">
						<span class="entry-plate">
							<PawPrint size={11} class="entry-mark" />
							<span class="entry-title">{entry.title}</span>
						</span>
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
					<span class="mat placeholder">
						<BambooSprout size={64} grow class="placeholder-sprout" />
						<span class="placeholder-label">Image coming soon</span>
					</span>
				{:else}
					<span class="mat">
						<img src={entry.image} alt={entry.title} onerror={() => markFailed(entry.id)} />
					</span>
				{/if}
			</span>
			<div class="lightbox-body">
				<span class="entry-plate">
					<PawPrint size={12} class="entry-mark" />
					<h2>{entry.title}</h2>
				</span>
				<p>{entry.caption}</p>
				{#if entry.url}
					<a class="lightbox-link" href={entry.url} target="_blank" rel="noreferrer">{linkLabel}</a>
				{/if}
			</div>
		</div>
	{/if}
</dialog>

<style>
	.gallery-page {
		max-width: var(--container-max);
		margin: 0 auto;
		padding: var(--space-4) var(--space-2) var(--space-6);
	}

	.gallery-header {
		max-width: 640px;
		margin: 0 auto;
		text-align: center;
	}

	.gallery-header h1 {
		font: var(--text-headline-lg-mobile);
		margin-bottom: var(--space-1);
	}

	.gallery-header p {
		margin: 0;
		color: var(--color-text-muted);
		font: var(--text-body-lg);
	}

	/* a fallen bamboo stalk laid under the header, standing in for a plain
	   hairline rule — reused as-is wherever Gallery is mounted */
	.grove-divider {
		display: flex;
		justify-content: center;
		margin: var(--space-3) 0 var(--space-4);
		color: var(--color-accent);
		opacity: 0.75;
	}

	.grove-divider :global(.grove-divider-stalk) {
		transform: rotate(90deg);
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
		outline-offset: 2px;
	}

	/* the "frame": a wood-toned mount hung on the gallery wall, with a small
	   nail head at the top edge — the mat sits inset inside it */
	.frame {
		position: relative;
		display: block;
		width: 100%;
		aspect-ratio: 4 / 5;
		padding: 0.5rem;
		border-radius: var(--radius);
		background: linear-gradient(
			155deg,
			color-mix(in srgb, var(--color-panda-rust) 32%, var(--color-surface-alt)) 0%,
			var(--color-surface-alt) 65%
		);
		border: 1px solid var(--color-border);
		transition:
			transform 0.35s ease,
			box-shadow 0.35s ease;
	}

	.frame::before {
		content: '';
		position: absolute;
		top: -0.3rem;
		left: 50%;
		translate: -50% 0;
		width: 0.4rem;
		height: 0.4rem;
		border-radius: var(--radius-full);
		background: color-mix(in srgb, var(--color-ink) 55%, var(--color-panda-rust) 45%);
		box-shadow: 0 1px 2px color-mix(in srgb, var(--color-ink) 35%, transparent);
	}

	.mat {
		position: relative;
		display: flex;
		width: 100%;
		height: 100%;
		border-radius: calc(var(--radius) - 0.2rem);
		overflow: hidden;
		background: var(--color-surface);
		border: 1px solid color-mix(in srgb, var(--color-ink) 10%, transparent);
	}

	.mat img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.3s ease;
	}

	.gallery-trigger:hover .mat img,
	.gallery-trigger:focus-visible .mat img {
		transform: scale(1.03);
	}

	.gallery-trigger:hover .frame,
	.gallery-trigger:focus-visible .frame {
		transform: translateY(-4px);
		box-shadow:
			0 12px 24px color-mix(in srgb, var(--color-ink) 22%, transparent),
			0 0 0 1px color-mix(in srgb, var(--color-accent) 25%, transparent);
	}

	:global([data-theme='bamboo']) .gallery-trigger:hover .frame,
	:global([data-theme='bamboo']) .gallery-trigger:focus-visible .frame {
		box-shadow:
			0 12px 24px color-mix(in srgb, var(--color-ink) 45%, transparent),
			0 0 1.5rem color-mix(in srgb, var(--color-accent) 45%, transparent);
	}

	/* leaves react to the frame's hover/focus, not just their own — a small
	   sign of life for the six placeholders sitting quietly in the grid */
	:global(.gallery-trigger:hover .placeholder-sprout .sprout-leaf-a),
	:global(.gallery-trigger:focus-visible .placeholder-sprout .sprout-leaf-a) {
		transform: rotate(-6deg);
	}

	:global(.gallery-trigger:hover .placeholder-sprout .sprout-leaf-b),
	:global(.gallery-trigger:focus-visible .placeholder-sprout .sprout-leaf-b) {
		transform: rotate(6deg);
	}

	.placeholder {
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-1);
		padding: var(--space-2);
		text-align: center;
		color: var(--color-text-muted);
	}

	.placeholder-label {
		font: var(--text-label-sm);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--color-accent-strong);
	}

	:global([data-theme='bamboo']) .placeholder-label {
		color: var(--color-canopy);
	}

	.entry-body {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.3rem;
		padding-top: var(--space-2);
		text-align: center;
	}

	.entry-plate {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.entry-plate :global(.entry-mark) {
		color: var(--color-accent);
		flex-shrink: 0;
	}

	.entry-title {
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.05rem;
		letter-spacing: 0.01em;
		color: var(--color-text);
	}

	.entry-caption {
		font: var(--text-label-sm);
		text-transform: uppercase;
		letter-spacing: 0.08em;
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
		transition: background 0.2s ease;
	}

	.lightbox-close:hover {
		background: color-mix(in srgb, var(--color-accent-strong) 65%, transparent);
	}

	.lightbox-frame {
		padding: var(--space-2);
		border-radius: var(--radius-lg) var(--radius-lg) 0 0;
		aspect-ratio: 4 / 5;
	}

	.lightbox-frame::before {
		top: 0.05rem;
		width: 0.55rem;
		height: 0.55rem;
	}

	.lightbox-frame .mat {
		border-radius: var(--radius);
	}

	.lightbox-body {
		padding: var(--space-3);
	}

	.lightbox-body .entry-plate {
		justify-content: flex-start;
		margin-bottom: 0.35rem;
	}

	.lightbox-body h2 {
		margin: 0;
	}

	.lightbox-body p {
		margin: 0;
		color: var(--color-text-muted);
	}

	.lightbox-link {
		display: inline-block;
		margin-top: var(--space-2);
		font: var(--text-body-md);
		color: var(--color-accent-strong);
	}

	.lightbox-link:hover {
		text-decoration: underline;
	}
</style>
