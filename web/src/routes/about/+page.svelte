<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import Chip from '$lib/components/Chip.svelte';
	import about from '$lib/data/about.json';
</script>

<svelte:head>
	<title>PandaLily — About</title>
</svelte:head>

<section class="about-page">
	<header class="page-header">
		<h1>About</h1>
		<p class="lede">
			A little lore about the spirit guardian of the bamboo forest — species, duty, likes, and
			dislikes.
		</p>
		<p class="editor-note">Placeholder bio below — real text coming soon.</p>
	</header>

	<div class="bento">
		<Card class="profile-card">
			<div class="profile-head">
				<div class="avatar" aria-hidden="true">PL</div>
				<div>
					<h2>{about.name}</h2>
					<p class="role">{about.role}</p>
				</div>
			</div>

			<div class="tags">
				{#each about.tags as tag (tag)}
					<Chip>{tag}</Chip>
				{/each}
			</div>

			<p class="bio">{about.bio}</p>

			<div class="duty">
				<h3>Duty</h3>
				<p>{about.duty}</p>
			</div>
		</Card>

		<Card class="vitals-card">
			<h3>Vitals</h3>
			<ul class="vitals-list">
				{#each about.vitals as vital (vital.label)}
					<li>
						<span class="vital-label">{vital.label}</span>
						<span class="vital-value">{vital.value}</span>
					</li>
				{/each}
			</ul>
		</Card>

		<Card class="preferences-card">
			<h3>Affinities &amp; Aversions</h3>
			<div class="preferences">
				<div>
					<h4>Affinities</h4>
					<ul class="pref-list affinities">
						{#each about.affinities as item (item)}
							<li>{item}</li>
						{/each}
					</ul>
				</div>
				<div>
					<h4>Aversions</h4>
					<ul class="pref-list aversions">
						{#each about.aversions as item (item)}
							<li>{item}</li>
						{/each}
					</ul>
				</div>
			</div>
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

	.editor-note {
		font: var(--text-label-sm);
		font-style: italic;
		color: var(--color-text-muted);
		opacity: 0.8;
		margin: 0;
	}

	.bento {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-3);
	}

	:global(.profile-card) {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	:global(.vitals-card),
	:global(.preferences-card) {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.profile-head {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.avatar {
		flex-shrink: 0;
		width: 4rem;
		height: 4rem;
		border-radius: 50%;
		background: var(--color-surface-alt);
		color: var(--color-accent-strong);
		border: 1px solid var(--color-border);
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-heading);
		font-weight: 700;
		font-size: 1.5rem;
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
		margin: 0;
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

	.bento h3 {
		font: var(--text-label);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-muted);
		margin: 0;
	}

	.bento h4 {
		font: var(--text-label-sm);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-muted);
		margin: 0 0 var(--space-1);
	}

	.duty p {
		font: var(--text-body-md);
		color: var(--color-text);
		margin: 0;
	}

	.vitals-list {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.vitals-list li {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: var(--space-2);
		padding: 0.6rem 0;
		border-bottom: 1px solid var(--color-border);
	}

	.vitals-list li:last-child {
		border-bottom: none;
	}

	.vital-label {
		font: var(--text-body-md);
		color: var(--color-text-muted);
	}

	.vital-value {
		font: var(--text-label);
		color: var(--color-ink);
		text-align: right;
	}

	.preferences {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-2);
	}

	.pref-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.pref-list li {
		position: relative;
		padding-left: 1rem;
		font: var(--text-body-md);
		color: var(--color-text);
	}

	.pref-list li::before {
		content: '';
		position: absolute;
		left: 0;
		top: 0.55em;
		width: 0.4rem;
		height: 0.4rem;
		border-radius: 50%;
	}

	.affinities li::before {
		background: var(--color-accent);
	}

	.aversions li::before {
		background: var(--color-text-muted);
	}

	@media (min-width: 768px) {
		h1 {
			font: var(--text-headline-lg);
		}

		.bento {
			grid-template-columns: 1fr 1fr;
		}
	}
</style>
