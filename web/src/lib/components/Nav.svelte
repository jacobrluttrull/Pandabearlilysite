<script lang="ts">
	import { page } from '$app/state';

	const tabs = [
		{ label: 'About', href: '/about' },
		{ label: 'Ref Sheet', href: '/refs' },
		{ label: 'Credits', href: '/credits' },
		{ label: 'Soundboard', href: '/soundboard' }
	];
</script>

<header class="nav">
	<div class="nav-inner">
		<a class="wordmark" href="/">PandaLily</a>
		<nav aria-label="Main">
			<ul>
				{#each tabs as tab}
					<li>
						<a href={tab.href} class:active={page.url.pathname === tab.href}>{tab.label}</a>
					</li>
				{/each}
			</ul>
		</nav>
	</div>
</header>

<style>
	.nav {
		position: sticky;
		top: 0;
		z-index: 50;
		background: color-mix(in srgb, var(--color-surface) 85%, transparent);
		backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--color-border);
	}

	.nav-inner {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		row-gap: 0.5rem;
		padding: 1rem 1.5rem;
	}

	.wordmark {
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.5rem;
		color: var(--color-text);
		text-decoration: none;
	}

	nav ul {
		display: flex;
		gap: 1.5rem;
		list-style: none;
		margin: 0;
		padding: 0;
	}

	nav a {
		text-decoration: none;
		color: var(--color-text-muted);
		font-family: var(--font-body);
		font-weight: 600;
		font-size: 0.9rem;
		padding: 0.25rem 0;
		border-bottom: 2px solid transparent;
		white-space: nowrap;
		transition: color 0.2s;
	}

	nav a:hover {
		color: var(--color-accent);
	}

	nav a.active {
		color: var(--color-accent);
		border-bottom-color: var(--color-accent);
	}

	@media (max-width: 640px) {
		nav ul {
			gap: 0.85rem;
		}

		nav a {
			font-size: 0.8rem;
		}

		.wordmark {
			font-size: 1.25rem;
		}
	}
</style>
