<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import PawPrint from './PawPrint.svelte';

	type CategoryIconProps = {
		/** which credits bucket this glyph represents */
		kind?: 'model' | 'assets' | 'emotes' | 'music';
		size?: number;
		class?: string;
	} & HTMLAttributes<SVGSVGElement>;

	let { kind = 'assets', size = 22, class: className = '', ...rest }: CategoryIconProps = $props();
</script>

<!-- small thematic glyphs so the four credit categories read as distinct
     buckets instead of visually-identical cards: a paw for the model
     itself, a folded screen for stream assets, a sunbeam spark for
     emotes, and a leaf-note for music. -->
{#if kind === 'model'}
	<PawPrint {size} class="category-icon-svg {className}" {...rest} />
{:else if kind === 'emotes'}
	<svg
		viewBox="0 0 24 24"
		width={size}
		height={size}
		class="category-icon-svg {className}"
		fill="currentColor"
		aria-hidden="true"
		{...rest}
	>
		<path
			d="M12 2 C12.8 7 14 9 19 10 C14 11 12.8 13 12 18 C11.2 13 10 11 5 10 C10 9 11.2 7 12 2 Z"
		/>
		<circle cx="19.2" cy="4" r="1.3" opacity="0.65" />
	</svg>
{:else if kind === 'music'}
	<svg
		viewBox="0 0 24 24"
		width={size}
		height={size}
		class="category-icon-svg {className}"
		fill="currentColor"
		aria-hidden="true"
		{...rest}
	>
		<ellipse cx="7.4" cy="18" rx="4.3" ry="3.1" transform="rotate(-15 7.4 18)" />
		<rect x="10.7" y="4" width="1.8" height="14.4" rx="0.9" />
		<path d="M12.5 4 C16.6 4 18.7 7 18.4 10.2 C15.3 9.1 12.5 7.4 12.5 5.8 Z" opacity="0.85" />
	</svg>
{:else}
	<svg
		viewBox="0 0 24 24"
		width={size}
		height={size}
		class="category-icon-svg {className}"
		fill="currentColor"
		aria-hidden="true"
		{...rest}
	>
		<rect x="3" y="4" width="17" height="13" rx="2.2" opacity="0.92" />
		<path d="M14 4 L20 4 L20 10 Z" opacity="0.55" />
		<rect x="9" y="18.6" width="6" height="1.6" rx="0.8" opacity="0.6" />
	</svg>
{/if}
