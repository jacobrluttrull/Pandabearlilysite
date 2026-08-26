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
	<header class="page-header">
		<h1>{about.name}</h1>
		<p class="lede">
			Field notes on the spirit guardian of the Zenith Bamboo Grove — species, duty, and a
			running log of what she loves (and loathes).
		</p>
	</header>

	<div class="registry-grid">
		<Card class="dossier-card" rail>
			<div class="dossier-head">
				<div class="id-block">
					<div class="specimen">
						<div class="avatar-ring">
							<div class="avatar" aria-hidden="true">PL</div>
						</div>
					</div>
					<div class="id-info">
						<h2>{about.name}</h2>
						<p class="role">{about.role}</p>
						<div class="tags">
							{#each about.tags as tag (tag)}
								<Chip>{tag}</Chip>
							{/each}
						</div>
					</div>
				</div>
			</div>

			<p class="bio">{about.bio}</p>

			<div class="duty">
				<h3>Duty</h3>
				<p>{about.duty}</p>
			</div>
		</Card>

		<Card class="vitals-card">
			<h3>Vitals</h3>
			<dl class="ledger">
				{#each about.vitals as vital (vital.label)}
					<div class="ledger-row">
						<dt>{vital.label}</dt>
						<span class="leader" aria-hidden="true"></span>
						<dd>{vital.value}</dd>
					</div>
				{/each}
			</dl>
		</Card>

		<Card class="provisions-card">
			<h3>Provisions</h3>
			<ul class="provisions-list">
				{#each about.affinities as item (item)}
					<li><PawPrint size={12} /> {item}</li>
				{/each}
			</ul>
		</Card>

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

	.page-header {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
		text-align: center;
		margin-bottom: var(--space-4);
	}

	h1 {
		font: var(--text-headline-lg-mobile);
		color: var(--color-ink);
		margin: 0;
	}

	.lede {
		font: var(--text-body-lg);
		color: var(--color-text-muted);
		max-width: 640px;
		margin: 0 auto;
	}

	.registry-grid {
		display: grid;
		grid-template-columns: 1fr;
		align-items: start;
		gap: var(--space-3);
	}

	:global(.dossier-card) {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	:global(.vitals-card),
	:global(.provisions-card),
	:global(.hazards-card) {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.dossier-head {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-3);
	}

	.id-block {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		flex: 1 1 260px;
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
		width: 4rem;
		height: 4rem;
		border-radius: 50%;
		background: color-mix(in srgb, var(--color-accent) 18%, var(--color-surface-alt) 82%);
		color: var(--color-accent-strong);
		border: 1px solid var(--color-border);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.5rem;
	}

	/* bamboo: keep the flat badge fill it already had — the green text is
	   unchanged, only cream gains the tinted background */
	:global([data-theme='bamboo']) .avatar {
		background: var(--color-surface-alt);
	}

	h2 {
		font: var(--text-headline-md);
		color: var(--color-ink);
		margin: 0;
	}

	.role {
		font: var(--text-label);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-accent-strong);
		margin: 0.15rem 0 var(--space-1);
	}

	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-1);
	}

	.bio {
		font: var(--text-body-md);
		color: var(--color-text);
		margin: 0;
	}

	.duty {
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

	.duty p {
		font: var(--text-body-md);
		color: var(--color-text);
		margin: 0;
	}

	/* Vitals: a field-guide spec ledger with dotted leaders between label and value */
	.ledger {
		display: flex;
		flex-direction: column;
		gap: 0.55rem;
		margin: 0;
		padding: 0;
	}

	.ledger-row {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		margin: 0;
	}

	.ledger-row dt {
		font: var(--text-body-md);
		color: var(--color-text-muted);
		white-space: nowrap;
	}

	.ledger-row .leader {
		flex: 1;
		min-width: 0.75rem;
		border-bottom: 2px dotted var(--color-border);
		margin-bottom: 0.3em;
	}

	.ledger-row dd {
		font: var(--text-label);
		color: var(--color-ink);
		margin: 0;
		white-space: nowrap;
		text-align: right;
	}

	/* Provisions: reuse the site's paw-print bullet motif */
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

	@media (min-width: 768px) {
		h1 {
			font: var(--text-headline-lg);
		}

		.registry-grid {
			grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
		}
	}
</style>
