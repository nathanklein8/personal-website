<script lang="ts">
  export let visible = false;
</script>

<header
  class="fixed top-0 left-0 w-full bg-white/80 backdrop-blur-lg transition-all duration-500 shadow-md"
  class:opacity-0={!visible}
  class:pointer-events-none={!visible}
>
  <nav class="flex justify-center gap-6 py-4 text-xl mx-6">
    <a class="font-bold" href="/">Nathan Klein</a>
    <span class="flex grow"></span>
    <a href="/photography">Photography</a>
    <a href="/hiking">Hiking</a>
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
