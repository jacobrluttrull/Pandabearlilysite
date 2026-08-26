<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import PawPrint from './PawPrint.svelte';

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
					{@const isActive = page.url.pathname === tab.href}
					<li>
						<a href={tab.href} class:active={isActive}>
							{tab.label}
							<PawPrint size={9} class="active-paw" aria-hidden="true" />
						</a>
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

	/* grid, not flex space-between: with a wide wordmark and a narrow toggle
	   button, space-between visually skews the centered tab list off to one
	   side. A 1fr/auto/1fr grid keeps the tabs genuinely centered regardless
	   of how the two side elements compare in width. */
	.nav-inner {
		max-width: 1200px;
		margin: 0 auto;
		padding: 1rem 1.5rem;
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		grid-template-areas: 'brand nav toggle';
		align-items: center;
		column-gap: 1rem;
	}

	.wordmark {
		grid-area: brand;
		justify-self: start;
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.5rem;
		color: var(--color-text);
		text-decoration: none;
	}

	nav {
		grid-area: nav;
		justify-self: center;
		min-width: 0;
	}

	nav ul {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 1.5rem;
		list-style: none;
		margin: 0;
		padding: 0;
	}

	nav a {
		position: relative;
		display: inline-flex;
		align-items: center;
		text-decoration: none;
		color: var(--color-text-muted);
		font-family: var(--font-body);
		font-weight: 600;
		font-size: 0.9rem;
		padding: 0.25rem 0 0.6rem;
		white-space: nowrap;
		transition: color 0.2s;
	}

	nav a:hover {
		color: var(--color-accent);
	}

	nav a.active {
		color: var(--color-accent);
	}

	nav a :global(.active-paw) {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%) scale(0.5);
		color: var(--color-accent);
		opacity: 0;
		transition:
			opacity 0.2s ease,
			transform 0.2s ease;
	}

	nav a.active :global(.active-paw) {
		opacity: 1;
		transform: translateX(-50%) scale(1);
	}

	.theme-toggle {
		grid-area: toggle;
		justify-self: end;
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

	/* below this width the single-row grid gets too cramped (tabs start
	   crowding the toggle button, then wrap raggedly). Reflow instead to a
	   predictable two-row layout: brand + toggle stay paired on row one,
	   the tab list gets its own full-width centered row below. */
	@media (max-width: 700px) {
		.nav-inner {
			grid-template-columns: 1fr auto;
			grid-template-areas:
				'brand toggle'
				'nav nav';
			row-gap: 0.65rem;
		}

		nav {
			justify-self: stretch;
		}

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
