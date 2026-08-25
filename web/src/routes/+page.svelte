<script lang="ts">
	import links from '$lib/data/links.json';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Chip from '$lib/components/Chip.svelte';
</script>

<svelte:head>
	<title>PandaLily</title>
</svelte:head>

<section class="hero">
	<div class="hero-avatar" aria-hidden="true">PL</div>

	<Chip class="eyebrow-chip">Welcome to the bamboo forest</Chip>

	<h1>PandaLily</h1>
	<p class="tagline">Spirit guardian of the bamboo forest</p>
	<p class="hero-sub">
		Clips, streams, and every other corner of the forest — all gathered in one quiet place.
	</p>

	<div class="hero-actions">
		<Button href="/soundboard">Enter the Soundboard</Button>
		<Button href="#connect" variant="secondary">Find Me Elsewhere</Button>
	</div>
</section>

<section class="directory" id="connect">
	<div class="directory-heading">
		<h2>Elsewhere in the Forest</h2>
		<p>Catch streams, chat, and everything else, right here.</p>
	</div>

	<a class="featured-link" href="/soundboard">
		<Card class="featured-card">
			<div class="featured-icon" aria-hidden="true">♪</div>
			<div class="featured-body">
				<Chip class="eyebrow-chip">Featured</Chip>
				<h3>Official Soundboard</h3>
				<p>Clips, alerts, and sounds from the bamboo forest.</p>
			</div>
			<span class="chevron" aria-hidden="true">→</span>
		</Card>
	</a>

	<ul class="link-list">
		{#each links as link (link.label)}
			<li>
				<a href={link.url} target="_blank" rel="noreferrer">
					<span class="icon" aria-hidden="true">{link.abbr}</span>
					<span class="label">{link.label}</span>
					<span class="chevron" aria-hidden="true">↗</span>
				</a>
			</li>
		{/each}
	</ul>
</section>

<style>
	/* hero */
	.hero {
		max-width: 640px;
		margin: 0 auto;
		padding: var(--space-6) var(--space-3) var(--space-4);
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-1);
	}

	.hero-avatar {
		width: 9rem;
		height: 9rem;
		border-radius: 50%;
		background: var(--color-surface-alt);
		color: var(--color-accent-strong);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-heading);
		font-size: 2.5rem;
		font-weight: 800;
		border: 1px solid var(--color-border);
		box-shadow:
			0 0 0 8px color-mix(in srgb, var(--color-accent) 12%, transparent),
			0 0 2.5rem color-mix(in srgb, var(--color-accent) 25%, transparent);
		margin-bottom: var(--space-2);
	}

	.hero :global(.eyebrow-chip) {
		color: var(--color-accent-strong);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		margin-bottom: var(--space-1);
	}

	.hero h1 {
		font: var(--text-headline-lg-mobile);
		margin: 0;
	}

	.hero .tagline {
		font: var(--text-body-lg);
		font-style: italic;
		color: var(--color-text-muted);
		margin: 0;
	}

	.hero-sub {
		font: var(--text-body-md);
		color: var(--color-text-muted);
		max-width: 34rem;
		margin: var(--space-1) 0 0;
	}

	.hero-actions {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: var(--space-2);
		margin-top: var(--space-3);
	}

	@media (min-width: 768px) {
		.hero h1 {
			font: var(--text-headline-lg);
		}
	}

	/* directory */
	.directory {
		max-width: 700px;
		margin: 0 auto;
		padding: 0 var(--space-3) var(--space-6);
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.directory-heading {
		text-align: center;
	}

	.directory-heading h2 {
		font: var(--text-headline-md);
		margin: 0 0 0.25rem;
	}

	.directory-heading p {
		margin: 0;
		color: var(--color-text-muted);
	}

	.featured-link {
		display: block;
		color: var(--color-text);
		text-decoration: none;
	}

	.featured-link :global(.featured-card) {
		position: relative;
		display: flex;
		align-items: center;
		gap: var(--space-3);
		overflow: hidden;
		border-width: 2px;
		border-color: var(--color-accent);
		transition: background-color 0.2s;
	}

	.featured-link :global(.featured-card)::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 4px;
		background: var(--color-accent);
	}

	.featured-link:hover :global(.featured-card) {
		background: var(--color-surface-alt);
	}

	.featured-icon {
		flex-shrink: 0;
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
		background: var(--color-surface-alt);
		color: var(--color-accent-strong);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
	}

	.featured-body {
		flex: 1;
	}

	.featured-body h3 {
		font: var(--text-headline-md);
		margin: 0 0 0.25rem;
	}

	.featured-body p {
		margin: 0;
		color: var(--color-text-muted);
		font: var(--text-body-md);
	}

	.chevron {
		color: var(--color-text-muted);
		font-size: 1.1rem;
	}

	.link-list {
		list-style: none;
		margin: 0;
		padding: 0;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	.link-list li + li {
		border-top: 1px solid var(--color-border);
	}

	.link-list a {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		text-decoration: none;
		color: var(--color-text);
		transition: background-color 0.2s;
	}

	.link-list a:hover {
		background: var(--color-surface-alt);
	}

	.icon {
		flex-shrink: 0;
		width: 2.25rem;
		height: 2.25rem;
		border-radius: 50%;
		background: var(--color-surface-alt);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-body);
		font-weight: 700;
		font-size: 0.7rem;
		color: var(--color-text-muted);
	}

	.label {
		flex: 1;
		font: var(--text-body-md);
		font-weight: 600;
	}
</style>
