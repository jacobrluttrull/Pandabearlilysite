<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	type CardProps = {
		variant?: 'default' | 'flat';
		class?: string;
		children: Snippet;
	} & HTMLAttributes<HTMLDivElement>;

	let { variant = 'default', class: className = '', children, ...rest }: CardProps = $props();
</script>

<div class="card {variant} {className}" {...rest}>
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
</style>
