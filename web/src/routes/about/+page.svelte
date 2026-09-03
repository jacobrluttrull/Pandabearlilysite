<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import Chip from '$lib/components/Chip.svelte';
	import PawPrint from '$lib/components/PawPrint.svelte';
	import about from '$lib/data/about.json';
</script>

{#snippet warningIcon()}
	<svg class="hazard-icon" viewBox="0 0 24 24" width="16" height="16" aria-hidden="true">
		<path d="M12 3 L22 20 H2 Z" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round" />
		<line x1="12" y1="9" x2="12" y2="14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
		<circle cx="12" cy="17" r="1" fill="currentColor" />
	</svg>
{/snippet}

<svelte:head>
	<title>PandaLily — About</title>
</svelte:head>

<section class="about-page">
	<div class="registry-grid">
		<Card class="dossier-card" rail>
			<div class="dossier-head">
				<div class="id-block">
					<div class="specimen">
						<div class="avatar-ring">
							<img class="avatar" src="/images/pfp.png" alt="PandaLily" />
						</div>
					</div>
					<div class="id-info">
						<h1>{about.name}</h1>
						<dl class="identity-facts">
							<div>
								<dt>Role</dt>
								<dd>{about.role}</dd>
							</div>
							<div>
								<dt>Species</dt>
								<dd>{about.species}</dd>
							</div>
							<div>
								<dt>Home</dt>
								<dd>{about.home}</dd>
							</div>
						</dl>
					</div>
				</div>
			</div>

			<div class="bio">
				<h3>Bio</h3>
				<p>{about.bio.join(' ')}</p>
			</div>

			<div class="lore">
				<h3>Lore</h3>
				<p>{about.lore.join(' ')}</p>
			</div>
		</Card>

		<Card class="favorite-things-card">
			<h3>Favorite Things</h3>
			<ul class="provisions-list">
				{#each about.affinities as item (item)}
					<li><PawPrint size={12} /> {item}</li>
				{/each}
			</ul>
		</Card>

		{#if about.youllLikeMeIf.length > 0}
			<Card class="compatibility-card">
				<h3>You’ll like me if...</h3>
				<ul class="compatibility-list">
					{#each about.youllLikeMeIf as interest (interest)}
						<li><Chip>{interest}</Chip></li>
					{/each}
				</ul>
			</Card>
		{/if}

		<Card class="hazards-card">
			<h3>Hazard Log</h3>
			<div class="stripe" aria-hidden="true"></div>
			<ul class="hazard-list">
				{#each about.aversions as item (item)}
					<li class="hazard-tile">{@render warningIcon()}<span>{item}</span></li>
				{/each}
			</ul>
		</Card>
	</div>
</section>

<style>
	.about-page {
		max-width: 1000px;
		margin: 0 auto;
		padding: var(--space-6) var(--space-3) var(--space-6);
	}

	h1 {
		font: var(--text-headline-lg-mobile);
		color: var(--color-ink);
		margin: 0;
	}

	.registry-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-3);
	}

	:global(.dossier-card) {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	:global(.favorite-things-card),
	:global(.compatibility-card),
	:global(.hazards-card) {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.dossier-head {
		display: flex;
		justify-content: center;
	}

	.id-block {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-3);
		margin-inline: auto;
	}

	.specimen {
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.35rem;
	}

	.avatar-ring {
		padding: 4px;
		border: 1px dashed var(--color-border);
		border-radius: 50%;
	}

	.avatar {
		width: 8rem;
		height: 8rem;
		border-radius: 50%;
		object-fit: cover;
		object-position: top center;
		border: 1px solid var(--color-border);
	}

	.identity-facts {
		display: grid;
		grid-template-columns: max-content minmax(0, 1fr);
		gap: 0.3rem var(--space-2);
		margin: var(--space-1) 0 0;
	}

	.identity-facts div {
		display: contents;
	}

	.identity-facts dt {
		font: var(--text-label);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-muted);
	}

	.identity-facts dd {
		font: var(--text-body-md);
		color: var(--color-ink);
		margin: 0;
	}

	.compatibility-list {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-1);
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.compatibility-list li {
		display: flex;
	}

	.bio {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.bio p {
		font: var(--text-body-md);
		color: var(--color-text);
		margin: 0;
	}

	.lore {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding-top: var(--space-2);
		border-top: 1px solid var(--color-border);
	}

	.registry-grid h3 {
		font: var(--text-label);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-muted);
		margin: 0;
	}

	.bio h3,
	.lore h3 {
		font: var(--text-heading-sm);
		text-transform: none;
		letter-spacing: 0.02em;
		color: var(--color-accent-strong);
	}

	.lore p {
		font: var(--text-body-md);
		color: var(--color-text);
		margin: 0;
	}

	/* Favorite Things: reuse the site's paw-print bullet motif. */
	.provisions-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.provisions-list li {
		display: flex;
		align-items: center;
		gap: 0.55rem;
		font: var(--text-body-md);
		color: var(--color-text);
	}

	.provisions-list :global(svg) {
		flex-shrink: 0;
		color: var(--color-accent);
	}

	/* Hazard Log: aversions get real visual weight — a caution-stripe header
	   and hand-stamped, slightly tilted warning tiles instead of a plain list.
	   This is the one place --color-panda-rust carries real meaning (danger),
	   so it earns being used a bit more boldly than elsewhere on the site. */
	:global(.hazards-card) h3 {
		color: var(--color-panda-rust);
	}

	.stripe {
		height: 5px;
		border-radius: var(--radius-full);
		margin: -0.35rem 0 0;
		background: repeating-linear-gradient(
			135deg,
			var(--color-panda-rust) 0 8px,
			var(--color-surface-alt) 8px 16px
		);
		opacity: 0.85;
	}

	.hazard-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.hazard-tile {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		padding: 0.55rem 0.75rem;
		border: 1px solid color-mix(in srgb, var(--color-panda-rust) 45%, var(--color-border) 55%);
		border-left: 3px solid var(--color-panda-rust);
		border-radius: var(--radius-sm);
		background: color-mix(in srgb, var(--color-panda-rust) 8%, var(--color-surface) 92%);
		font: var(--text-body-md);
		color: var(--color-text);
		transition:
			transform 0.15s ease,
			box-shadow 0.15s ease;
	}

	.hazard-tile:nth-child(odd) {
		transform: rotate(-0.6deg);
	}

	.hazard-tile:nth-child(even) {
		transform: rotate(0.6deg);
	}

	.hazard-tile:hover {
		transform: translateY(-2px) rotate(0deg);
		box-shadow: 0 6px 16px color-mix(in srgb, var(--color-panda-rust) 25%, transparent);
	}

	.hazard-icon {
		flex-shrink: 0;
		color: var(--color-panda-rust);
	}

	@media (max-width: 479px) {
		.id-block {
			flex-direction: column;
		}

		.id-info {
			width: min(100%, 18rem);
		}

		h1 {
			text-align: center;
		}
	}

	@media (min-width: 768px) {
		h1 {
			font: var(--text-headline-lg);
		}

		.registry-grid {
			grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
		}
	}
</style>
