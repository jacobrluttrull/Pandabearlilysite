<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	type CardProps = {
		variant?: 'default' | 'flat';
		/** Adds the "Accent Rail" treatment — a 3px forest-green top edge for
		 *  featured/prominent cards (cream only; bamboo keeps its normal border). */
		rail?: boolean;
		class?: string;
		children: Snippet;
	} & HTMLAttributes<HTMLDivElement>;

	let {
		variant = 'default',
		rail = false,
		class: className = '',
		children,
		...rest
	}: CardProps = $props();
</script>

<div class="card {variant} {className}" class:rail {...rest}>
	{@render children()}
</div>

<style>
	.card {
		border-radius: var(--radius-lg);
		padding: var(--space-3);
		border: 1px solid var(--surface-glass-border);
		background: var(--surface-glass-bg);
		backdrop-filter: var(--surface-blur);
		-webkit-backdrop-filter: var(--surface-blur);
	}

	/* flat variant: always an opaque solid surface, no blur, for content
	   nested inside another card or anywhere the glass treatment is too much */
	.card.flat {
		background: var(--color-surface);
		backdrop-filter: none;
		-webkit-backdrop-filter: none;
	}

	/* Accent Rail (cream): a 3px forest-green top edge for featured cards —
	   "a pop of color without overwhelming the beige palette" per the design
	   doc. Bamboo already reads as accent-forward via its glass treatment, so
	   it uses the normal glass-border color. The width stays 3px in both themes
	   so switching themes cannot shift the card's contents or surrounding grid. */
	.card.rail {
		border-top: 3px solid var(--color-accent);
	}

	:global([data-theme='bamboo']) .card.rail {
		border-top-color: var(--surface-glass-border);
	}
</style>
