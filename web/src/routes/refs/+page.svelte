<script lang="ts">
	import Gallery from '$lib/components/Gallery.svelte';
	import refSheet from '$lib/data/ref-sheet.json';
</script>

<svelte:head>
	<title>PandaLily — Ref Sheet</title>
</svelte:head>

<Gallery
	heading="Reference Sheet"
	subtitle="Official visual reference for PandaLily — body, expressions, colors, and outfits."
	entries={refSheet}
/>
