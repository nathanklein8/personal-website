<script lang="ts">
	import { enhance } from '$app/forms';
	import { buttonVariants } from '$lib/components/ui/button';
	import { Card, CardContent } from '$lib/components/ui/card';
	import ErrorCard from '$lib/components/error-card.svelte';
	import type { PageProps } from './$types';
	import { Check, X, Camera, Trash2, Plus } from '@lucide/svelte';

	let { data, form }: PageProps = $props();

	// simple state to toggle "Add New" form
	let showAddForm = $state(false);
</script>

<section id="photos-editor" class="flex flex-col items-center gap-4 py-5 mt-16">
	<div class="flex w-full max-w-4xl items-center justify-between px-3">
		<h1 class="font-code text-lg">Manage Photo Gallery</h1>
		<button 
			class={buttonVariants({ size: 'sm' })} 
			onclick={() => showAddForm = !showAddForm}
		>
			<Plus class="mr-2 w-4 h-4" /> {showAddForm ? 'Cancel' : 'Add Photo'}
		</button>
	</div>

	{#if showAddForm}
		<Card class="mx-3 w-2xl shadow-lg border-primary/20">
			<CardContent class="p-6">
				<form
					class="flex flex-col gap-4"
					method="POST"
					action="?/addPhoto"
					use:enhance
				>
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<label class="flex flex-col gap-1">
							Title
							<input required name="title" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							File Path (relative to /photos)
							<input required name="filePath" placeholder="landscapes/sunset.jpg" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Alt Text
							<input name="altText" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Location
							<input name="location" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Date Taken
							<input name="dateTaken" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Camera
							<input name="camera" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Lens
							<input name="lens" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Aperture (e.g. 2.8)
							<input name="aperture" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Shutter Speed (e.g. 1/500)
							<input name="shutterSpeed" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							ISO
							<input name="iso" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex flex-col gap-1">
							Sort Order
							<input type="number" name="sortOrder" value="0" class="rounded-md border p-2 text-sm text-muted-foreground" />
						</label>
						<label class="flex items-center gap-2 mt-6">
							<input type="checkbox" name="visible" checked /> Visible in gallery
						</label>
					</div>
					<div class="flex items-center gap-2">
						<button class={buttonVariants({ size: 'sm' })}> Save Photo </button>
						{#if form?.success && !form?.id}
							<Check color={'green'} />
						{:else if form?.failure && !form?.id}
							<X color={'red'} /> {form.message}
						{/if}
					</div>
				</form>
			</CardContent>
		</Card>
	{/if}

	<div class="flex flex-col gap-4 w-full max-w-4xl px-3">
		{#each data.photos as photo}
			<Card class="shadow-md overflow-hidden">
				<CardContent class="p-0 flex flex-col sm:flex-row">
					<div class="w-full sm:w-48 h-48 bg-muted flex items-center justify-center overflow-hidden">
						<img 
							src={`${data.apiURL}/api/photos/${photo.id}/image`} 
							alt={photo.altText || photo.title}
							class="w-full h-full object-cover"
						/>
					</div>
					<div class="p-4 flex-grow">
						<form
							method="POST"
							action="?/updatePhoto"
							use:enhance
							class="flex flex-col gap-3"
						>
							<input type="hidden" name="id" value={photo.id} />
							<div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">TITLE</span>
									<input name="title" value={photo.title} required class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">PATH</span>
									<input name="filePath" value={photo.filePath} required class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">SORT</span>
									<input type="number" name="sortOrder" value={photo.sortOrder} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">CAMERA</span>
									<input name="camera" value={photo.camera || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">LENS</span>
									<input name="lens" value={photo.lens || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">ISO</span>
									<input name="iso" value={photo.iso || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">APERTURE</span>
									<input name="aperture" value={photo.aperture || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">SHUTTER</span>
									<input name="shutterSpeed" value={photo.shutterSpeed || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
								<label class="flex flex-col gap-1">
									<span class="text-xs font-bold opacity-50">LOCATION</span>
									<input name="location" value={photo.location || ''} class="rounded border p-1 text-sm text-muted-foreground" />
								</label>
							</div>

							<div class="flex items-center justify-between mt-2">
								<div class="flex items-center gap-2">
									<input type="checkbox" name="visible" checked={photo.visible} />
									<span class="text-xs text-muted-foreground">Visible</span>
								</div>
								
								<div class="flex items-center gap-2">
									<button class={buttonVariants({ size: 'sm', variant: 'ghost' })}> Update </button>
									{#if form?.success && form?.id == photo.id}
										<Check color={'green'} class="w-4 h-4" />
									{:else if form?.failure && form?.id == photo.id}
										<X color={'red'} class="w-4 h-4" />
									{/if}
								</div>
							</div>
						</form>
					</div>
					<div class="w-12 flex items-center justify-center">
						<form method="POST" action="?/deletePhoto" use:enhance>
							<input type="hidden" name="id" value={photo.id} />
							<button class={buttonVariants({ variant: 'destructive', size: 'icon' })}>
								<Trash2 class="w-4 h-4" />
							</button>
						</form>
					</div>
				</CardContent>
			</Card>
		{/each}
	</div>
</section>
