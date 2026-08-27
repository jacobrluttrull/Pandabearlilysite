<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import Chip from '$lib/components/Chip.svelte';
	import CategoryIcon from '$lib/components/CategoryIcon.svelte';
	import credits from '$lib/data/credits.json';

	type CategoryMeta = {
		slug: string;
		icon: 'model' | 'assets' | 'emotes' | 'music';
		color: string;
	};

	// ties each bucket to a glyph and one accent from the grove palette, used
	// once each, so the four cards read as distinct places rather than
	// identical boxes with different labels.
	const CATEGORY_META: Record<string, CategoryMeta> = {
		'Live2D Model': { slug: 'model', icon: 'model', color: 'var(--color-accent)' },
		'Stream Assets & Screens': { slug: 'assets', icon: 'assets', color: 'var(--color-ink)' },
		Emotes: { slug: 'emotes', icon: 'emotes', color: 'var(--color-canopy-strong)' },
		Music: { slug: 'music', icon: 'music', color: 'var(--color-panda-rust)' }
	};
	const FALLBACK_META: CategoryMeta = { slug: 'other', icon: 'assets', color: 'var(--color-accent)' };

	// a role with no name on file yet — real and accurate, not a bug.
	function isPending(name: string): boolean {
		return !name?.trim() || name.trim().toUpperCase() === 'TBD';
	}

	// tracks which entries' <img> failed to load, keyed by category+role, so
	// we can drop the thumbnail instead of showing a broken-image glyph
	let failed = $state<Record<string, boolean>>({});

	function markFailed(key: string) {
		failed[key] = true;
	}
</script>

<svelte:head>
	<title>PandaLily — Credits</title>
</svelte:head>

<section class="credits-page">
	<header class="page-header">
		<h1>Credits</h1>
	</header>

	<div class="category-grid">
		{#each credits as { category, entries } (category)}
			{@const meta = CATEGORY_META[category] ?? FALLBACK_META}
			<Card class="category-card {meta.slug}" rail>
				<div class="category-head">
					<span class="category-icon" style="color: {meta.color}">
						<CategoryIcon kind={meta.icon} size={24} />
					</span>
					<h2>{category}</h2>
				</div>
				<ul class="entry-list">
					{#each entries as { role, name, url, image } (role)}
						{@const pending = isPending(name)}
						{@const key = `${category}-${role}`}
						<li class="entry" class:pending>
							{#if image && !failed[key]}
								<img
									class="entry-thumb"
									src={image}
									alt=""
									loading="lazy"
									onerror={() => markFailed(key)}
								/>
							{/if}
							<div class="entry-info">
								<Chip>{role}</Chip>
								{#if pending}
									<!-- no name on file yet — the image and role speak for themselves -->
								{:else if url}
									<a class="entry-name" href={url} target="_blank" rel="noreferrer">{name}</a>
								{:else}
									<span class="entry-name">{name}</span>
								{/if}
							</div>
						</li>
					{/each}
				</ul>

				{#if category === 'Live2D Model'}
					<div class="showcase">
						<h3 class="showcase-heading">Model Showcase</h3>
						<div class="showcase-frame">
							<iframe
								src="https://www.youtube.com/embed/RkBsFK242Kk"
								title="PandaLily Live2D Model Showcase"
								loading="lazy"
								allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
								referrerpolicy="strict-origin-when-cross-origin"
								allowfullscreen
							></iframe>
						</div>
					</div>
				{/if}
			</Card>
		{/each}
	</div>
</section>

<style>
	.credits-page {
		max-width: var(--container-max);
		margin: 0 auto;
		padding: var(--space-6) var(--space-2) var(--space-6);
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.page-header {
		text-align: center;
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.page-header h1 {
		font: var(--text-headline-lg-mobile);
		margin: 0;
		color: var(--color-text);
	}

	.category-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-3);
	}

	:global(.category-card) {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.category-head {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.category-icon {
		display: inline-flex;
		flex-shrink: 0;
	}

	:global(.category-card) h2 {
		margin: 0;
		font: var(--text-headline-md);
		color: var(--color-text);
	}

	/* rail color echoes each category's glyph accent (cream only — bamboo's
	   glass card already reads accent-forward, so it keeps the plain border
	   from Card's own bamboo override rather than stacking a second cue) */
	:global([data-theme='cream'] .category-card.rail.assets) {
		border-top-color: var(--color-ink);
	}

	:global([data-theme='cream'] .category-card.rail.emotes) {
		border-top-color: var(--color-canopy-strong);
	}

	:global([data-theme='cream'] .category-card.rail.music) {
		border-top-color: var(--color-panda-rust);
	}

	/* each credit is its own big, obviously-visible plate — a grid of tiles
	   instead of a cramped row list, since a thumbnail small enough to sit
	   inline with text reads as decoration, not as the actual asset */
	.entry-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(15rem, 1fr));
		gap: var(--space-3);
	}

	.entry {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-2);
		padding: var(--space-2);
		border-radius: var(--radius-lg);
		border: 1px solid var(--color-border);
		background: var(--color-surface-alt);
	}

	.entry-thumb {
		width: 100%;
		aspect-ratio: 1 / 1;
		max-width: 15rem;
		padding: 0.75rem;
		border-radius: var(--radius);
		object-fit: contain;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
	}

	.entry-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.35rem;
		min-width: 0;
	}

	/* Live2D Model's card gets a bonus block below its entry list — a
	   16:9 embed of the model showcase, framed like the rest of the card's
	   content rather than bleeding to the card's edges */
	.showcase {
		margin-top: var(--space-3);
		padding-top: var(--space-3);
		border-top: 1px solid var(--color-border);
	}

	:global([data-theme='bamboo']) .showcase {
		border-top-color: color-mix(in srgb, var(--color-ink) 20%, transparent);
	}

	.showcase-heading {
		margin: 0 0 var(--space-2);
		font: var(--text-label-sm);
		text-transform: uppercase;
		letter-spacing: 0.06em;
		color: var(--color-text-muted);
	}

	.showcase-frame {
		position: relative;
		width: 100%;
		aspect-ratio: 16 / 9;
		border-radius: var(--radius);
		overflow: hidden;
		border: 1px solid var(--color-border);
		background: var(--color-surface-alt);
	}

	.showcase-frame iframe {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		border: 0;
	}

	.entry-name {
		font: var(--text-body-lg);
		font-weight: 700;
		color: var(--color-text);
		text-decoration: none;
	}

	a.entry-name:hover {
		color: var(--color-accent-strong);
		text-decoration: underline;
	}

	@media (min-width: 768px) {
		.page-header h1 {
			font: var(--text-headline-lg);
		}

		.category-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}
</style>
