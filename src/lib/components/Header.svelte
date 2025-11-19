<script lang="ts">
	import ModeToggle from './mode-toggle.svelte';
	import {
		Sheet,
		SheetContent,
		SheetTrigger,
		SheetHeader,
		SheetClose
	} from '$lib/components/ui/sheet';
	import { buttonVariants } from '$lib/components/ui/button';
	import { Menu } from '@lucide/svelte';

	export let visible = false;
</script>

<header
	class="font-code fixed top-0 left-0 flex h-12 w-full flex-col justify-center bg-white/70 shadow-lg backdrop-blur-xl transition-all duration-500 sm:h-16 dark:bg-black/70"
	class:opacity-0={!visible}
	class:pointer-events-none={!visible}
>
	<nav class="text-lg flex flex-row items-center justify-center gap-3 px-4 sm:px-6">
		<a href="/" class="font-semibold hover:font-bold">Nathan Klein</a>

		<span class="flex grow"></span>

		<!-- Desktop links -->
		<div class="hidden items-center gap-4 sm:flex">
			<a href="/#projects" class="hover:font-semibold">Projects</a>
			<a href="/photography" class="hover:font-semibold">Photography</a>
			<a href="/hiking" class="hover:font-semibold">Hiking</a>
			<ModeToggle />
		</div>

		<!-- Mobile menu -->
		<div class="flex items-center gap-2 sm:hidden">
			<ModeToggle />
			<Sheet>
				<SheetTrigger>
					<div class={buttonVariants({variant: "ghost", size: "icon-sm"})}>
						<Menu />
					</div>
				</SheetTrigger>
				<SheetContent side="right" class="w-2/5 min-w-fit">
					<SheetHeader class="font-code h-12 border-b"></SheetHeader>
					<div class="mx-4 flex flex-col gap-4 text-base items-end text-lg">
						<SheetClose>
							<a href="/" class="hover:underline">About Me</a>
						</SheetClose>
						<SheetClose>
							<a href="/#projects" class="hover:underline">Projects</a>
						</SheetClose>
						<SheetClose>
							<a href="/photography" class="hover:underline">Photography</a>
						</SheetClose>
						<SheetClose>
							<a href="/hiking" class="hover:underline">Hiking</a>
						</SheetClose>
					</div>
				</SheetContent>
			</Sheet>
		</div>
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
