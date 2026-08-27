<script lang="ts">
	import Gallery from '$lib/components/Gallery.svelte';
	import artists from '$lib/data/art.json';
</script>

<svelte:head>
	<title>PandaLily — Art</title>
</svelte:head>

<Gallery
	heading="Art"
	subtitle="Collection of Pandalily artwork, I thank all the artists for their work and support!"
	entries={artists}
	linkLabel="Visit portfolio"
/>
