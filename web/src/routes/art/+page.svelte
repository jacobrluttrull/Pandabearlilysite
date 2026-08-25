<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import Chip from '$lib/components/Chip.svelte';
	import artists from '$lib/data/art.json';
</script>

<svelte:head>
	<title>PandaLily — Art</title>
</svelte:head>

<section class="art-page">
	<header class="page-header">
		<h1>Art</h1>
		<p>Commissioned art made by these wonderful artists — credited with thanks.</p>
	</header>

	<Card class="artist-card">
		<ul class="entry-list">
			{#each artists as { name, role, url } (name)}
				<li class="entry">
					<Chip>{role}</Chip>
					{#if url}
							<a class="entry-name" href={url} target="_blank" rel="noreferrer">{name}</a>
						{:else}
							<span class="entry-name">{name}</span>
						{/if}
				</li>
			{/each}
		</ul>
	</Card>
</section>

<style>
	.art-page {
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

	.page-header p {
		margin: 0;
		font: var(--text-body-lg);
		color: var(--color-text-muted);
	}

	:global(.artist-card) {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
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

	.entry-name:hover {
		color: var(--color-accent-strong);
		text-decoration: underline;
	}

	@media (min-width: 768px) {
		.page-header h1 {
			font: var(--text-headline-lg);
		}
	}
</style>
