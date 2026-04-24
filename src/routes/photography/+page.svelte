<script lang="ts">
	import PhotoInfo from '$lib/components/photo-info.svelte';
	import { Card } from '@nk/shared/components/ui/card';
	let { data } = $props();

	let selectedPhoto = $state<any>(null);
	let isDialogOpen = $state(false);

	function openPhoto(photo: any) {
		selectedPhoto = photo;
		isDialogOpen = true;
	}
</script>

<section class="flex flex-col items-center py-16 px-4 sm:px-8">
	<h1 class="text-4xl font-bold mb-8 font-code">Photography</h1>

	{#if data.photos.length === 0}
		<div class="text-center py-20">
			<p class="text-muted-foreground">No photos in the gallery yet.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 w-full max-w-7xl">
			{#each data.photos as photo}
				<Card 
					class="group relative cursor-pointer overflow-hidden rounded-xl transition-all hover:ring-2 ring-primary"
					onclick={() => openPhoto(photo)}
				>
					<div class="aspect-square overflow-hidden bg-muted">
						<img 
							src={`/api/photos/${photo.year}/${photo.event}/${photo.filename}/image`} 
							alt={photo.altText || photo.title}
							class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
							loading="lazy"
						/>
					</div>
					<div class="absolute bottom-0 left-0 right-0 p-3 bg-linear-to-t from-black/80 to-transparent opacity-0 group-hover:opacity-100 transition-opacity">
						<p class="text-white text-sm font-medium truncate">{photo.title}</p>
					</div>
				</Card>
			{/each}
		</div>
	{/if}
	<PhotoInfo 
		open={isDialogOpen} 
		onOpenChange={(open) => {
			isDialogOpen = open;
			if (!open) selectedPhoto = null;
		}} 
		photo={{
			...selectedPhoto,
			src: selectedPhoto ? `/api/photos/${selectedPhoto.year}/${selectedPhoto.event}/${selectedPhoto.filename}/image` : '',
			alt: selectedPhoto?.altText,
			date: selectedPhoto?.dateTaken,
			shutter: selectedPhoto?.shutterSpeed
		}}
	/>
</section>
