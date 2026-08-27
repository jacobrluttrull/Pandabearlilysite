<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import BambooSprout from '$lib/components/BambooSprout.svelte';
	import refSheet from '$lib/data/ref-sheet.json';

	type RefEntry = {
		id: string;
		title: string;
		caption: string;
		image: string;
	};

	const entries = refSheet as RefEntry[];

	function pad(n: number) {
		return String(n).padStart(2, '0');
	}

	// tracks which entries' <img> failed to load, so we fall back to the
	// "reference pending" stamp instead of a broken-image glyph
	let failed = $state<Record<string, boolean>>({});

	function markFailed(id: string) {
		failed[id] = true;
	}

	function showsImage(entry: RefEntry) {
		return entry.image && !failed[entry.id];
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

	<div class="plates">
		{#each entries as entry, i (entry.id)}
			<Card class="plate-card" rail>
				<div class="plate-index">
					<span class="plate-num">{pad(i + 1)}<span class="plate-total">/{pad(entries.length)}</span></span>
					<span class="plate-leader" aria-hidden="true"></span>
					<h2 class="plate-title">{entry.title}</h2>
				</div>

				<div class="plate-window">
					<span class="crop tl" aria-hidden="true"></span>
					<span class="crop tr" aria-hidden="true"></span>
					<span class="crop bl" aria-hidden="true"></span>
					<span class="crop br" aria-hidden="true"></span>

					{#if showsImage(entry)}
						<img src={entry.image} alt={entry.title} onerror={() => markFailed(entry.id)} />
					{:else}
						<div class="plate-stamp">
							<BambooSprout size={30} class="plate-sprout" />
							<span class="stamp-text">Reference<br />Pending</span>
						</div>
					{/if}
				</div>

				<p class="plate-caption">
					<span class="fig-tag">FIG.{pad(i + 1)}</span>
					{entry.caption}
				</p>
			</Card>
		{/each}
	</div>
</section>

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

	.plates {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-3);
	}

	@media (min-width: 760px) {
		.plates {
			/* deliberately caps at two columns (not three, like Art's gallery
			   wall) — fewer, larger plates read as "study this closely",
			   not "browse many thumbnails" */
			grid-template-columns: repeat(2, 1fr);
		}
	}

	.plates :global(.plate-card) {
		display: flex;
		flex-direction: column;
	}

	.plate-index {
		display: flex;
		align-items: baseline;
		gap: 0.6rem;
		margin-bottom: var(--space-2);
	}

	.plate-num {
		flex-shrink: 0;
		font: var(--text-mono-md);
		color: var(--color-accent-strong);
		letter-spacing: 0.02em;
	}

	:global([data-theme='bamboo']) .plate-num {
		color: var(--color-canopy);
	}

	.plate-total {
		color: var(--color-text-muted);
		font-weight: 600;
	}

	.plate-leader {
		flex: 1 1 auto;
		min-width: 0.75rem;
		height: 0;
		border-bottom: 1px dashed var(--color-border);
	}

	.plate-title {
		flex-shrink: 0;
		margin: 0;
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 0.03em;
		color: var(--color-text);
		text-align: right;
	}

	/* the "plate": a technical drawing sheet — graph-paper fill, print-style
	   crop marks at each corner — replacing Gallery's wood-toned photo frame */
	.plate-window {
		position: relative;
		width: 100%;
		aspect-ratio: 16 / 10;
		border-radius: var(--radius-sm);
		overflow: hidden;
		border: 1px solid var(--color-border);
		background-color: var(--color-surface-alt);
		background-image:
			linear-gradient(color-mix(in srgb, var(--color-border) 55%, transparent) 1px, transparent 1px),
			linear-gradient(90deg, color-mix(in srgb, var(--color-border) 55%, transparent) 1px, transparent 1px);
		background-size: 18px 18px;
	}

	.plate-window img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.crop {
		position: absolute;
		width: 0.55rem;
		height: 0.55rem;
		border-color: var(--color-accent);
		opacity: 0.65;
		pointer-events: none;
	}

	.crop.tl {
		top: 0.25rem;
		left: 0.25rem;
		border-top: 2px solid;
		border-left: 2px solid;
	}

	.crop.tr {
		top: 0.25rem;
		right: 0.25rem;
		border-top: 2px solid;
		border-right: 2px solid;
	}

	.crop.bl {
		bottom: 0.25rem;
		left: 0.25rem;
		border-bottom: 2px solid;
		border-left: 2px solid;
	}

	.crop.br {
		bottom: 0.25rem;
		right: 0.25rem;
		border-bottom: 2px solid;
		border-right: 2px solid;
	}

	/* "reference pending" — a rotated rubber-stamp treatment, distinct from
	   Gallery's nursery placard, to read as an official document marker */
	.plate-stamp {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.3rem;
		text-align: center;
	}

	.plate-stamp::before {
		content: '';
		position: absolute;
		inset: 18% 14%;
		border: 2px dashed color-mix(in srgb, var(--color-panda-rust) 65%, transparent);
		border-radius: var(--radius-sm);
		rotate: -5deg;
	}

	.plate-stamp :global(.plate-sprout) {
		position: relative;
		rotate: -5deg;
	}

	.stamp-text {
		position: relative;
		rotate: -5deg;
		font: var(--text-mono-sm);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		line-height: 1.35;
		color: color-mix(in srgb, var(--color-panda-rust) 80%, var(--color-text));
	}

	.plate-caption {
		margin: var(--space-2) 0 0;
		font: var(--text-body-md);
		color: var(--color-text-muted);
	}

	.fig-tag {
		display: inline-block;
		margin-right: 0.4em;
		font: var(--text-mono-sm);
		letter-spacing: 0.04em;
		color: var(--color-accent-strong);
	}

	:global([data-theme='bamboo']) .fig-tag {
		color: var(--color-canopy);
	}
</style>
