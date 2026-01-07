<script lang="ts">
	import ModeToggle from '$lib/components/mode-toggle.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { Menu } from '@lucide/svelte';

	export let visible = false;
</script>

<header
	class="font-code fixed top-0 z-1 left-0 flex h-12 w-full flex-col justify-center bg-white/70 shadow-lg backdrop-blur-xl transition-all duration-500 sm:h-16 dark:bg-black/70"
	class:opacity-0={!visible}
	class:pointer-events-none={!visible}
>
	<nav class="text-lg flex flex-row items-center justify-center gap-3 px-4 sm:px-6">
		<a href="/" class="font-semibold hover:font-bold">Edit Site Contents</a>
		<span class="flex grow"></span>
		<ModeToggle />
	</nav>
</header>

<style>
	header {
		opacity: 1;
	}
	.opacity-0 {
		opacity: 0;
	}
</style>
