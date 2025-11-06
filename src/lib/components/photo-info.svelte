<script lang="ts">
	import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "$lib/components/ui/dialog";
	import { Separator } from "$lib/components/ui/separator";
	import { ScrollArea } from "$lib/components/ui/scroll-area";
	import { X } from "@lucide/svelte";
	import { Button } from "$lib/components/ui/button";

	export let open = false;
	export let onOpenChange: (value: boolean) => void = () => {};

	export let photo = {
		title: "Untitled",
		src: "",
		alt: "",
		date: "",
		location: "",
		camera: "",
		lens: "",
		aperture: "",
		shutter: "",
		iso: ""
	};
</script>

<Dialog {open} {onOpenChange}>
	<DialogContent class="max-w-3xl p-0 overflow-hidden">
		<!-- Image Section -->
		<div class="relative w-full bg-black">
			<img
				src={photo.src}
				alt={photo.alt || photo.title}
				class="w-full max-h-[70vh] object-contain select-none"
				loading="lazy"
			/>
			<Button
				variant="ghost"
				size="icon"
				class="absolute top-2 right-2 text-white hover:bg-white/20"
				onclick={() => onOpenChange(false)}
			>
				<X class="w-5 h-5" />
			</Button>
		</div>

		<!-- Metadata Section -->
		<ScrollArea class="p-4 max-h-[40vh]">
			<DialogHeader>
				<DialogTitle class="text-lg font-semibold">{photo.title}</DialogTitle>
				<DialogDescription class="text-sm text-muted-foreground">
					{photo.location} {photo.date ? `• ${photo.date}` : ""}
				</DialogDescription>
			</DialogHeader>

			<Separator class="my-3" />

			<div class="grid grid-cols-2 sm:grid-cols-3 gap-4 text-sm">
				{#if photo.camera}
					<div>
						<span class="font-medium">Camera:</span>
						<div>{photo.camera}</div>
					</div>
				{/if}
				{#if photo.lens}
					<div>
						<span class="font-medium">Lens:</span>
						<div>{photo.lens}</div>
					</div>
				{/if}
				{#if photo.aperture}
					<div>
						<span class="font-medium">Aperture:</span>
						<div>ƒ/{photo.aperture}</div>
					</div>
				{/if}
				{#if photo.shutter}
					<div>
						<span class="font-medium">Shutter:</span>
						<div>{photo.shutter}s</div>
					</div>
				{/if}
				{#if photo.iso}
					<div>
						<span class="font-medium">ISO:</span>
						<div>{photo.iso}</div>
					</div>
				{/if}
			</div>
		</ScrollArea>
	</DialogContent>
</Dialog>

<style>
	img {
		user-select: none;
	}
</style>
