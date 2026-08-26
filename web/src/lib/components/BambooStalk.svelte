<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	type BambooStalkProps = {
		/** number of jointed segments — taller dividers want more */
		segments?: number;
		/** adds a small leaf sprig near the top */
		leaf?: boolean;
		/** applies the slow ambient sway loop (respects prefers-reduced-motion) */
		sway?: boolean;
		class?: string;
	} & HTMLAttributes<SVGSVGElement>;

	let { segments = 3, leaf = true, sway = false, class: className = '', ...rest }: BambooStalkProps =
		$props();

	const segmentHeight = 46;
	const gap = 6;
	const stalkWidth = 14;
	const stalkX = 4;
	// extra horizontal room so the leaf has space to sit fully inside the
	// viewBox instead of being clipped past its right edge
	const leafSpace = 16;

	const width = $derived(stalkX * 2 + stalkWidth + (leaf ? leafSpace : 0));
	const height = $derived(segments * segmentHeight + (segments - 1) * gap);
</script>

<!-- signature mark: a jointed bamboo stalk, reused as a section divider,
     card accent, and ambient background motif instead of a flat rule -->
<svg
	viewBox="0 0 {width} {height}"
	width={width}
	{height}
	class="bamboo-stalk {className}"
	class:sway
	fill="none"
	aria-hidden="true"
	{...rest}
>
	{#each Array(segments) as _, i}
		{@const y = i * (segmentHeight + gap)}
		<rect
			x={stalkX}
			y={y}
			width={stalkWidth}
			height={segmentHeight}
			rx={stalkWidth / 2}
			fill="currentColor"
			opacity={0.9 - i * 0.1}
		/>
		<!-- a soft vertical highlight for a rounded, cylindrical feel -->
		<rect
			x={stalkX + 2.5}
			y={y + 4}
			width="2"
			height={segmentHeight - 8}
			rx="1"
			fill="var(--color-bg)"
			opacity="0.16"
		/>
		{#if i > 0}
			<rect
				x={stalkX - 2}
				y={y - gap / 2 - 1.5}
				width={stalkWidth + 4}
				height="3"
				rx="1.5"
				fill="currentColor"
			/>
		{/if}
	{/each}

	{#if leaf}
		<path
			d="M {stalkX + stalkWidth - 3} 13
			   C {stalkX + stalkWidth + 6} 11, {width - 4} 6, {width - 1} 1
			   C {width - 3} 8, {stalkX + stalkWidth + 3} 12, {stalkX + stalkWidth - 5} 17 Z"
			fill="var(--color-canopy)"
		/>
	{/if}
</svg>

<style>
	.sway {
		transform-origin: 50% 100%;
		animation: bamboo-sway 6s ease-in-out infinite;
	}

	@keyframes bamboo-sway {
		0%,
		100% {
			transform: rotate(-1.5deg);
		}
		50% {
			transform: rotate(1.5deg);
		}
	}
</style>
