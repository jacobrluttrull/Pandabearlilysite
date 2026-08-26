<script lang="ts">
	import links from '$lib/data/links.json';
	import Chip from '$lib/components/Chip.svelte';
	import SocialIcon from '$lib/components/SocialIcon.svelte';
</script>

<svelte:head>
	<title>PandaLily</title>
</svelte:head>

<section class="hero">
	<div class="hero-avatar-wrap">
		<div class="canopy-glow" aria-hidden="true"></div>
		<div class="hero-avatar" aria-hidden="true">
			<span class="avatar-mark">PL</span>
		</div>
	</div>

	<Chip class="eyebrow-chip">Welcome to the bamboo forest</Chip>

	<h1>PandaLily</h1>
	<p class="tagline">Spirit guardian of the bamboo forest</p>
</section>

<section class="directory">
	<ul class="link-list">
		<li>
			<a href="/soundboard">
				<span class="icon" aria-hidden="true">
					<svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor" aria-hidden="true">
						<circle cx="6" cy="18" r="3.4" />
						<rect x="8.8" y="3" width="2" height="15.4" rx="1" />
						<path
							d="M10.8 3 C15.5 3.6 18 6.8 17.6 10.4 C14.3 9.4 10.8 7.4 10.8 5.4 Z"
							opacity="0.9"
						/>
					</svg>
				</span>
				<span class="label">Soundboard</span>
				<span class="chevron" aria-hidden="true">→</span>
			</a>
		</li>
		{#each links as link (link.label)}
			<li>
				<a href={link.url} target="_blank" rel="noreferrer">
					<span class="icon" aria-hidden="true"><SocialIcon platform={link.icon} size={17} /></span>
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
		padding: var(--space-4) var(--space-3);
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: var(--space-1);
	}

	/* staggered emergence: each hero element rises in in sequence, like
	   stepping out from under the canopy into a sunbeam. `both` keeps the
	   pre-animation state applied before the animation starts and the
	   final state applied after it ends. Covered by the global
	   prefers-reduced-motion safety net (collapses to a near-instant cut). */
	.hero-avatar-wrap,
	.hero :global(.eyebrow-chip),
	.hero h1,
	.hero .tagline {
		animation: rise-in 0.7s ease-out both;
	}

	.hero-avatar-wrap {
		animation-delay: 0s;
	}
	.hero :global(.eyebrow-chip) {
		animation-delay: 0.1s;
	}
	.hero h1 {
		animation-delay: 0.2s;
	}
	.hero .tagline {
		animation-delay: 0.3s;
	}

	@keyframes rise-in {
		from {
			opacity: 0;
			transform: translateY(14px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.hero-avatar-wrap {
		position: relative;
		flex-shrink: 0;
		width: 8.5rem;
		height: 8.5rem;
		margin-bottom: var(--space-2);
	}

	/* dappled sunbeam breaking through the canopy, behind the medallion */
	.canopy-glow {
		position: absolute;
		inset: -30%;
		z-index: 0;
		border-radius: 50%;
		background: radial-gradient(
			circle,
			color-mix(in srgb, var(--color-canopy) 55%, transparent) 0%,
			color-mix(in srgb, var(--color-canopy) 18%, transparent) 45%,
			transparent 72%
		);
		filter: blur(2px);
	}

	.hero-avatar {
		position: relative;
		z-index: 1;
		width: 100%;
		height: 100%;
		border-radius: 50%;
		background: var(--color-surface-alt);
		display: flex;
		align-items: center;
		justify-content: center;
		border: 2px dashed color-mix(in srgb, var(--color-canopy) 70%, var(--color-accent) 30%);
		box-shadow:
			inset 0 0 0 6px var(--color-bg),
			inset 0 0 0 8px color-mix(in srgb, var(--color-accent) 45%, transparent),
			0 0 2.5rem color-mix(in srgb, var(--color-canopy) 35%, transparent);
	}

	.avatar-mark {
		font-family: var(--font-heading);
		font-size: 2.25rem;
		font-weight: 800;
		color: var(--color-accent-strong);
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
		gap: var(--space-2);
	}

	.chevron {
		color: var(--color-text-muted);
		font-size: 1.1rem;
		transition: transform 0.2s ease;
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

	.link-list li {
		position: relative;
	}

	.link-list li + li {
		border-top: 1px solid var(--color-border);
	}

	/* accent rail: a bar that grows in on hover, echoing Card's "Accent
	   Rail" treatment used for featured content elsewhere on the site */
	.link-list li::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		bottom: 0;
		width: 3px;
		background: var(--color-accent);
		transform: scaleY(0);
		transition: transform 0.2s ease;
	}

	.link-list li:hover::before {
		transform: scaleY(1);
	}

	.link-list a {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3);
		text-decoration: none;
		color: var(--color-text);
		transition:
			background-color 0.2s,
			transform 0.2s ease;
	}

	.link-list a:hover {
		background: color-mix(in srgb, var(--color-accent) 10%, var(--color-surface-alt) 90%);
		transform: translateX(4px);
	}

	.link-list a:hover .chevron {
		transform: translateX(2px);
		color: var(--color-accent-strong);
	}

	/* bamboo: keep the plain alt-surface hover it already had */
	:global([data-theme='bamboo']) .link-list a:hover {
		background: var(--color-surface-alt);
	}

	/* leaf badges: an asymmetric, slightly rotated silhouette instead of a
	   plain monogram circle, cycling through the grove's three accent
	   colors (forest green / canopy gold / panda rust) so the directory
	   doesn't read as one repeated component */
	.icon {
		flex-shrink: 0;
		width: 2.4rem;
		height: 2.4rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 62% 38% 55% 45% / 45% 55% 38% 62%;
		transform: rotate(-6deg);
		font-family: var(--font-body);
		font-weight: 700;
		font-size: 0.7rem;
		transition: transform 0.2s ease;
	}

	.link-list li:nth-child(2n) .icon {
		transform: rotate(5deg);
	}

	.link-list a:hover .icon {
		transform: rotate(0deg) scale(1.06);
	}

	.link-list li:nth-child(3n + 1) .icon {
		background: color-mix(in srgb, var(--color-accent) 22%, var(--color-surface-alt) 78%);
		color: var(--color-accent-strong);
	}

	.link-list li:nth-child(3n + 2) .icon {
		background: color-mix(in srgb, var(--color-canopy) 30%, var(--color-surface-alt) 70%);
		color: var(--color-canopy-strong);
	}

	.link-list li:nth-child(3n) .icon {
		background: color-mix(in srgb, var(--color-panda-rust) 24%, var(--color-surface-alt) 76%);
		color: var(--color-panda-rust);
	}

	.label {
		flex: 1;
		font: var(--text-body-md);
		font-weight: 600;
	}
</style>
