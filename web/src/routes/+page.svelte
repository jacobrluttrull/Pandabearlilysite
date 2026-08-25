<script lang="ts">
	import links from '$lib/data/links.json';
	import Button from '$lib/components/Button.svelte';
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

	<div class="hero-actions">
		<Button href="/soundboard">Enter the Soundboard</Button>
	</div>
</section>

<section class="directory">
	<ul class="link-list">
		<li>
			<a href="/soundboard">
				<span class="icon" aria-hidden="true">♪</span>
				<span class="label">Soundboard</span>
				<span class="chevron" aria-hidden="true">→</span>
			</a>
		</li>
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
		padding: var(--space-4) var(--space-3) var(--space-3);
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
		font: var(--text-headline-md);
		color: var(--color-text);
		margin: 0;
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
		padding: 0 var(--space-3) var(--space-4);
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
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
		background: color-mix(in srgb, var(--color-accent) 10%, var(--color-surface-alt) 90%);
	}

	/* bamboo: keep the plain alt-surface hover it already had */
	:global([data-theme='bamboo']) .link-list a:hover {
		background: var(--color-surface-alt);
	}

	.icon {
		flex-shrink: 0;
		width: 2.25rem;
		height: 2.25rem;
		border-radius: 50%;
		background: color-mix(in srgb, var(--color-accent) 18%, var(--color-surface-alt) 82%);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-body);
		font-weight: 700;
		font-size: 0.7rem;
		color: var(--color-accent-strong);
	}

	/* bamboo: keep the muted flat badge it already had */
	:global([data-theme='bamboo']) .icon {
		background: var(--color-surface-alt);
		color: var(--color-text-muted);
	}

	.label {
		flex: 1;
		font: var(--text-body-md);
		font-weight: 600;
	}
</style>
