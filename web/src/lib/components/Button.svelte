<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes, HTMLAnchorAttributes } from 'svelte/elements';

	type ButtonProps = {
		variant?: 'primary' | 'secondary';
		size?: 'md' | 'sm';
		href?: string;
		class?: string;
		children: Snippet;
	} & HTMLButtonAttributes &
		HTMLAnchorAttributes;

	let {
		variant = 'primary',
		size = 'md',
		href,
		class: className = '',
		children,
		...rest
	}: ButtonProps = $props();
</script>

{#if href}
	<a {href} class="btn {variant} {size} {className}" {...rest}>
		{@render children()}
	</a>
{:else}
	<button class="btn {variant} {size} {className}" {...rest}>
		{@render children()}
	</button>
{/if}

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5em;
		border-radius: var(--radius-full);
		text-decoration: none;
		cursor: pointer;
		border: 1px solid transparent;
		transition:
			transform 0.15s ease,
			background-color 0.2s ease,
			border-color 0.2s ease,
			color 0.2s ease,
			box-shadow 0.2s ease;
	}

	.btn:hover {
		transform: translateY(-2px) scale(1.02);
	}

	.btn:active {
		transform: translateY(0) scale(0.99);
	}

	.btn:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	.btn:disabled {
		cursor: not-allowed;
		opacity: 0.5;
		transform: none;
	}

	/* sizes */
	.btn.md {
		padding: var(--space-2) var(--space-3);
		font: var(--text-label);
		letter-spacing: 0.02em;
	}

	.btn.sm {
		padding: var(--space-1) var(--space-2);
		font: var(--text-label-sm);
		letter-spacing: 0.02em;
	}

	/* variants (cream default) */
	.btn.primary {
		background: var(--color-ink);
		color: var(--color-on-ink);
		border-color: transparent;
	}

	.btn.secondary {
		background: transparent;
		color: var(--color-accent);
		border-color: var(--color-accent);
	}

	/* bamboo divergence: primary fills with accent instead of ink,
	   secondary gains a spirit-glow hover treatment */
	:global([data-theme='bamboo']) .btn.primary {
		background: var(--color-accent);
		color: var(--color-on-ink);
	}

	:global([data-theme='bamboo']) .btn.secondary:hover:not(:disabled) {
		border-color: var(--color-accent);
		box-shadow: 0 0 1.25rem color-mix(in srgb, var(--color-accent) 55%, transparent);
	}
</style>
