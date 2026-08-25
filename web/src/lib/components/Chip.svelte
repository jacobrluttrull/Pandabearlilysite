<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	type ChipProps = {
		class?: string;
		children: Snippet;
	} & HTMLAttributes<HTMLSpanElement>;

	let { class: className = '', children, ...rest }: ChipProps = $props();
</script>

<span class="chip {className}" {...rest}>
	{@render children()}
</span>

<style>
	.chip {
		display: inline-flex;
		align-items: center;
		gap: 0.35em;
		padding: calc(var(--space-1) / 2) var(--space-1);
		border-radius: var(--radius-full);
		border: 1px solid transparent;
		font: var(--text-label-sm);
		white-space: nowrap;

		/* cream default: a solid, quiet organizational tool */
		background: var(--color-surface-alt);
		color: var(--color-ink);
	}

	/* bamboo divergence: glassmorphic "leaf" treatment instead of a solid fill */
	:global([data-theme='bamboo']) .chip {
		background: var(--surface-glass-bg);
		border-color: var(--surface-glass-border);
		backdrop-filter: var(--surface-blur);
		-webkit-backdrop-filter: var(--surface-blur);
		color: var(--color-text);
	}
</style>
