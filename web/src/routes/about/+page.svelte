<script lang="ts">
	import links from '$lib/data/links.json';
</script>

<svelte:head>
	<title>PandaLily — Links</title>
</svelte:head>

<section class="links-page">
	<aside class="profile">
		<div class="avatar" aria-hidden="true">PL</div>
		<h1>PandaLily</h1>
		<p class="tagline">Spirit guardian of the bamboo forest</p>
	</aside>

	<div class="directory">
		<a class="featured" href="/soundboard">
			<div class="featured-icon" aria-hidden="true">♪</div>
			<div class="featured-body">
				<span class="eyebrow">Featured</span>
				<h2>Official Soundboard</h2>
				<p>Clips, alerts, and sounds from the bamboo forest.</p>
			</div>
			<span class="chevron" aria-hidden="true">→</span>
		</a>

		<ul class="link-list">
			{#each links as link}
				<li>
					<a href={link.url} target="_blank" rel="noreferrer">
						<span class="icon" aria-hidden="true">{link.abbr}</span>
						<span class="label">{link.label}</span>
						<span class="chevron" aria-hidden="true">↗</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
</section>

<style>
	.links-page {
		max-width: 1000px;
		margin: 0 auto;
		padding: 3rem 1.5rem 4rem;
		display: grid;
		grid-template-columns: 1fr;
		gap: 2.5rem;
	}

	.profile {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 0.5rem;
	}

	.avatar {
		width: 7rem;
		height: 7rem;
		border-radius: 50%;
		background: var(--color-surface-alt);
		color: var(--color-accent-strong);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-heading);
		font-size: 2rem;
		font-weight: 700;
		border: 1px solid var(--color-border);
		margin-bottom: 0.5rem;
	}

	h1 {
		font-size: 2rem;
		margin: 0;
	}

	.tagline {
		margin: 0;
		color: var(--color-text-muted);
	}

	.directory {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.featured {
		position: relative;
		display: flex;
		align-items: center;
		gap: 1.25rem;
		padding: 1.5rem;
		background: var(--color-surface);
		border: 2px solid var(--color-accent);
		border-radius: var(--radius-lg);
		text-decoration: none;
		color: var(--color-text);
		overflow: hidden;
		transition: background-color 0.2s;
	}

	.featured::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 4px;
		background: var(--color-accent);
	}

	.featured:hover {
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

	.eyebrow {
		display: inline-block;
		font-family: var(--font-body);
		font-size: 0.7rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-accent-strong);
		background: var(--color-surface-alt);
		padding: 0.15rem 0.5rem;
		border-radius: var(--radius-sm);
		margin-bottom: 0.5rem;
	}

	.featured h2 {
		font-size: 1.25rem;
		margin: 0 0 0.25rem;
	}

	.featured p {
		margin: 0;
		color: var(--color-text-muted);
		font-size: 0.9rem;
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
		gap: 1rem;
		padding: 1.1rem 1.25rem;
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
		font-family: var(--font-body);
		font-weight: 600;
	}

	@media (min-width: 768px) {
		.links-page {
			grid-template-columns: 260px 1fr;
			align-items: start;
		}

		.profile {
			position: sticky;
			top: 6rem;
			align-items: flex-start;
			text-align: left;
		}
	}
</style>
