<script lang="ts">
  import Header from '$lib/components/Header.svelte';
  import { onMount } from 'svelte';
  let showHeader = false;

  onMount(() => {
    const landing = document.getElementById('landing');
    const observer = new IntersectionObserver(
      ([entry]) => {
        showHeader = !entry.isIntersecting;
      },
      { threshold: 0.1 }
    );
    if (landing) observer.observe(landing);
    return () => observer.disconnect();
  });
</script>

<style>
  .font-pixel {
    font-family: 'Jersey 10', monospace;
    letter-spacing: 2px;
  }
</style>

<Header visible={showHeader} />

<section
	id="landing"
	class="flex h-screen flex-col items-center justify-center bg-gradient-to-b from-blue-500 to-indigo-600 text-white"
>
	<div class="pb-16 text-center">
		<h1 class="mb-4 text-5xl font-pixel">Hi, I’m Nathan</h1>
		<p class="text-lg text-background">Check out some of my projects!</p>
	</div>

	<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24"
		><path
			fill="currentColor"
			d="M11 4h2v12h2v2h-2v2h-2v-2H9v-2h2zM7 14v2h2v-2zm0 0v-2H5v2zm10 0v2h-2v-2zm0 0v-2h2v2z"
		/></svg
	>
</section>

<section class="flex h-[200vh] flex-col items-center justify-center bg-gray-50">
	<p class="text-xl text-gray-700 mb-8">Under Construction</p>
</section>
