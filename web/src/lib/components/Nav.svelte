<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	const tabs = [
		{ label: 'About', href: '/about' },
		{ label: 'Art', href: '/art' },
		{ label: 'Ref Sheet', href: '/refs' },
		{ label: 'Credits', href: '/credits' },
		{ label: 'Soundboard', href: '/soundboard' }
	];

	let theme = $state<'bamboo' | 'cream'>('bamboo');

	onMount(() => {
		theme = document.documentElement.getAttribute('data-theme') === 'cream' ? 'cream' : 'bamboo';
	});

	function toggleTheme() {
		theme = theme === 'bamboo' ? 'cream' : 'bamboo';
		document.documentElement.setAttribute('data-theme', theme);
		try {
			localStorage.setItem('theme', theme);
		} catch {}
	}
</script>

<header class="nav">
	<div class="nav-inner">
		<a class="wordmark" href="/">PandaLily</a>
		<nav aria-label="Main">
			<ul>
				{#each tabs as tab (tab.href)}
					<li>
						<a href={tab.href} class:active={page.url.pathname === tab.href}>{tab.label}</a>
					</li>
				{/each}
			</ul>
		</nav>
		<button
			type="button"
			class="theme-toggle"
			onclick={toggleTheme}
			aria-label={theme === 'bamboo' ? 'Switch to light theme' : 'Switch to dark theme'}
		>
			{#if theme === 'bamboo'}
				<svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
					<circle cx="12" cy="12" r="4.5" fill="currentColor" />
					<g stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
						<line x1="12" y1="1.5" x2="12" y2="4" />
						<line x1="12" y1="20" x2="12" y2="22.5" />
						<line x1="1.5" y1="12" x2="4" y2="12" />
						<line x1="20" y1="12" x2="22.5" y2="12" />
						<line x1="4.5" y1="4.5" x2="6.2" y2="6.2" />
						<line x1="17.8" y1="17.8" x2="19.5" y2="19.5" />
						<line x1="4.5" y1="19.5" x2="6.2" y2="17.8" />
						<line x1="17.8" y1="6.2" x2="19.5" y2="4.5" />
					</g>
				</svg>
			{:else}
				<svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
					<path
						fill="currentColor"
						d="M20.5 14.6a8.5 8.5 0 0 1-11.1-11 8.5 8.5 0 1 0 11.1 11Z"
					/>
				</svg>
			{/if}
		</button>
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

	.theme-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.25rem;
		height: 2.25rem;
		border-radius: var(--radius-full);
		border: 1px solid var(--color-border);
		background: transparent;
		color: var(--color-text-muted);
		cursor: pointer;
		transition:
			color 0.2s,
			border-color 0.2s;
	}

	.theme-toggle:hover {
		color: var(--color-accent);
		border-color: var(--color-accent);
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
