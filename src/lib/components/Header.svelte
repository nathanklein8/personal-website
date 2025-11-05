<script lang="ts">
	import ModeToggle from "./mode-toggle.svelte";

	export let visible = false;
</script>

<header
	class="font-code fixed top-0 left-0 flex h-12 sm:h-16 w-full flex-col justify-center bg-white/70 dark:bg-black/70 shadow-lg backdrop-blur-xl transition-all duration-500"
	class:opacity-0={!visible}
	class:pointer-events-none={!visible}
>
	<nav class="flex flex-row justify-center items-center gap-3 px-6 text-md md:text-lg">
		<a href="/">Nathan Klein</a>
		<span class="flex grow"></span>
		<a href="#projects">Projects</a>
		<!-- <a href="/photography">Photography</a>
		<a href="/hiking">Hikes</a> -->
		<ModeToggle></ModeToggle>
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
