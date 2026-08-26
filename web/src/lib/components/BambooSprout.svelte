<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	type BambooSproutProps = {
		size?: number;
		/** applies a slow ambient growth loop on the leaves (respects prefers-reduced-motion) */
		grow?: boolean;
		class?: string;
	} & HTMLAttributes<SVGSVGElement>;

	let { size = 40, grow = false, class: className = '', ...rest }: BambooSproutProps = $props();
</script>

<!-- Gallery-specific empty-slot mark: a bamboo sprout breaking through soil,
     standing in for "this tribute is still taking root" instead of a
     generic broken-image glyph. Leaves react to a `.gallery-trigger:hover`
     ancestor (wired up in Gallery.svelte) and/or the `grow` prop below. -->
<svg
	viewBox="0 0 32 36"
	width={size}
	height={size * (36 / 32)}
	class="bamboo-sprout {className}"
	class:grow
	aria-hidden="true"
	{...rest}
>
	<ellipse class="sprout-soil" cx="16" cy="32.5" rx="11" ry="2.6" fill="var(--color-border)" opacity="0.55" />
	<path
		class="sprout-stem"
		d="M16 32 C 15 24, 15.5 17, 16.5 11"
		stroke="var(--color-accent)"
		stroke-width="2.25"
		stroke-linecap="round"
		fill="none"
	/>
	<path
		class="sprout-leaf sprout-leaf-a"
		d="M16 20 C 9 19, 4.5 14.5, 5 8.5 C 11.5 9.5, 15.5 13.5, 16 20 Z"
		fill="var(--color-accent)"
	/>
	<path
		class="sprout-leaf sprout-leaf-b"
		d="M16.5 13 C 22 11, 26.5 6, 26 1 C 20 2.5, 16.5 6.5, 16.5 13 Z"
		fill="var(--color-canopy)"
	/>
	<circle class="sprout-bud" cx="16.5" cy="10.5" r="1.6" fill="var(--color-canopy-strong)" />
</svg>

<style>
	.sprout-leaf,
	.sprout-bud {
		transform-box: fill-box;
		transform-origin: bottom center;
		transition:
			transform 0.5s ease,
			opacity 0.5s ease;
	}

	/* ambient loop, used for the single large lightbox sprout only — grid
	   tiles stay calm and react to hover instead (see Gallery.svelte) */
	.grow .sprout-leaf-a {
		animation: sprout-sway-a 5s ease-in-out infinite;
	}

	.grow .sprout-leaf-b {
		animation: sprout-sway-b 5.5s ease-in-out infinite;
	}

	@keyframes sprout-sway-a {
		0%,
		100% {
			transform: rotate(-2deg);
		}
		50% {
			transform: rotate(3deg);
		}
	}

	@keyframes sprout-sway-b {
		0%,
		100% {
			transform: rotate(2deg);
		}
		50% {
			transform: rotate(-3deg);
		}
	}
</style>
