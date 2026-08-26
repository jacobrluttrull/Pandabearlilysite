<script lang="ts">
	import PawPrint from './PawPrint.svelte';
</script>

<footer class="site-footer">
	<PawPrint size={12} class="footer-paw" />
	<p>PandaLily — Spirit Guardian of the Bamboo Forest</p>
</footer>

<style>
	.site-footer {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.4rem;
		text-align: center;
		padding: 2rem 1rem;
		color: var(--color-text-muted);
		font-family: var(--font-body);
		font-size: 0.8rem;
	}

	.site-footer :global(.footer-paw) {
		color: var(--color-panda-rust);
		opacity: 0.7;
	}

	.site-footer p {
		margin: 0;
	}
</style>
