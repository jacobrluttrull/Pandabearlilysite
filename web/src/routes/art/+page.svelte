<script lang="ts">
	import Gallery from '$lib/components/Gallery.svelte';
	import artists from '$lib/data/art.json';
</script>

<svelte:head>
	<title>PandaLily — Art</title>
</svelte:head>

<Gallery
	heading="Art"
	subtitle="Commissioned art made by these wonderful artists — credited with thanks."
	entries={artists}
	linkLabel="Visit portfolio"
/>
