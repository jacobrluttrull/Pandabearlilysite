<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	type PawPrintProps = {
		size?: number;
		class?: string;
	} & HTMLAttributes<SVGSVGElement>;

	let { size = 14, class: className = '', ...rest }: PawPrintProps = $props();
</script>

<!-- signature mark: PandaLily's paw print, used as an active-state and
     list-bullet glyph across the site instead of a generic dot or bar -->
<svg
	viewBox="0 0 24 24"
	width={size}
	height={size}
	class="paw-print {className}"
	fill="currentColor"
	aria-hidden="true"
	{...rest}
>
	<ellipse cx="12" cy="16.2" rx="6.1" ry="5.1" />
	<ellipse cx="4.6" cy="10.4" rx="2.15" ry="2.7" transform="rotate(-18 4.6 10.4)" />
	<ellipse cx="9.6" cy="5.6" rx="2.15" ry="2.85" transform="rotate(-6 9.6 5.6)" />
	<ellipse cx="14.7" cy="5.5" rx="2.15" ry="2.85" transform="rotate(6 14.7 5.5)" />
	<ellipse cx="19.5" cy="10.2" rx="2.15" ry="2.7" transform="rotate(18 19.5 10.2)" />
</svg>
