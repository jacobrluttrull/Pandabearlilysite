<script lang="ts">
	import Button from './Button.svelte';

	type MediaActionsProps = {
		asset: string;
		label: string;
	};

	let { asset, label }: MediaActionsProps = $props();
</script>

<div class="media-actions" aria-label={`Actions for ${label}`}>
	<Button href={asset} target="_blank" rel="noreferrer" variant="secondary" size="sm">Open</Button>
	<Button href={asset} download variant="secondary" size="sm">Download</Button>
</div>

<style>
	.media-actions {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: var(--space-1);
		margin-top: var(--space-2);
	}
</style>
