<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from "mode-watcher";
	import Header from '$lib/components/Header.svelte';
	

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link href="https://fonts.googleapis.com/css2?family=Jersey+10&display=swap" rel="stylesheet" />
</svelte:head>

<Header visible={true} />

<ModeWatcher />

<div class="pt-16">
	{@render children?.()}
</div>
