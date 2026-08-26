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
</script>

<svelte:head>
	<title>PandaLily — Credits</title>
</svelte:head>

<section class="credits-page">
	<header class="page-header">
		<h1>Credits</h1>
		<p class="lede">
			Every artist, rigger, and composer who helped PandaLily take shape — plus a few seats
			still being set at the table for names yet to come.
		</p>
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
					{#each entries as { role, name, url } (role)}
						{@const pending = isPending(name)}
						<li class="entry" class:pending>
							<Chip>{role}</Chip>
							{#if pending}
								<span class="entry-name pending-name" title="Not yet credited">TBD</span>
							{:else if url}
								<a class="entry-name" href={url} target="_blank" rel="noreferrer">{name}</a>
							{:else}
								<span class="entry-name">{name}</span>
							{/if}
						</li>
					{/each}
				</ul>
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

	.page-header .lede {
		margin: 0;
		font: var(--text-body-lg);
		color: var(--color-text-muted);
		max-width: 640px;
		align-self: center;
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

	.entry-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
	}

	.entry {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) 0;
	}

	/* cream: thin solid "paper shadow" divider between list items */
	.entry + .entry {
		border-top: 1px solid var(--color-border);
	}

	/* bamboo: organic divider — low-opacity line that fades at the edges
	   instead of a hard rule spanning the full width */
	:global([data-theme='bamboo']) .entry + .entry {
		border-top: none;
		position: relative;
	}

	:global([data-theme='bamboo']) .entry + .entry::before {
		content: '';
		position: absolute;
		top: 0;
		left: 5%;
		right: 5%;
		height: 1px;
		background: linear-gradient(
			to right,
			transparent,
			color-mix(in srgb, var(--color-ink) 20%, transparent),
			transparent
		);
	}

	.entry-name {
		font: var(--text-body-md);
		color: var(--color-text);
		text-decoration: none;
	}

	a.entry-name:hover {
		color: var(--color-accent-strong);
		text-decoration: underline;
	}

	/* a quiet, dashed tag for roles not yet credited — reads as "pending",
	   not broken */
	.pending-name {
		display: inline-flex;
		align-items: center;
		padding: 0.15em 0.65em;
		border-radius: var(--radius-full);
		border: 1px dashed var(--color-border);
		font: var(--text-label-sm);
		letter-spacing: 0.04em;
		color: var(--color-text-muted);
	}

	@media (min-width: 768px) {
		.page-header h1 {
			font: var(--text-headline-lg);
		}

		.category-grid {
			grid-template-columns: repeat(2, 1fr);
			align-items: start;
		}
	}
</style>
