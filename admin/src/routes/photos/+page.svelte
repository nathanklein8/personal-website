<script lang="ts">
	import { enhance } from '$app/forms';
	import { buttonVariants } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$lib/components/ui/card';
	import ErrorCard from '@nk/shared/components/error-card';
	import type { PageProps } from './$types';
	import { Check, X, Trash2, Plus, ChevronRight, ChevronDown, RefreshCw, Camera, Folder, Image as ImageIcon } from '@lucide/svelte';

	let { data, form }: PageProps = $props();

	// File explorer state
	let years = $state<string[]>([]);
	let selectedYear = $state<string | null>(null);
	let events = $state<string[]>([]);
	let selectedEvent = $state<string | null>(null);
	let photos = $state<string[]>([]);
	let selectedPhoto = $state<string | null>(null);
	let previewUrl = $state<string | null>(null);
	let loading = $state(false);
	let explorerError = $state<string | null>(null);

	// Add photo state
	let showAddForm = $state(false);
	let addTitle = $state('');
	let addSortOrder = $state(0);
	let addSuccess = $state<number | null>(null);
	let addFailure = $state<string | null>(null);

	// Regenerate thumbnails
	let regenerating = $state(false);
	let regenerateSuccess = $state(false);
	let regenerateFailure = $state<string | null>(null);

	// Expand/collapse state for years
	let expandedYears = $state<Set<string>>(new Set());

	// Expand/collapse state for events
	let expandedEvents = $state<Set<string>>(new Set());

	// Load available years
	async function loadYears() {
		loading = true;
		explorerError = null;
		try {
			const res = await fetch(`${data.apiURL}/api/photos/available`);
			if (!res.ok) throw new Error(`Failed to load years: ${res.status}`);
			years = await res.json();
		} catch (e: any) {
			explorerError = e.message;
		} finally {
			loading = false;
		}
	}

	// Load events for a year
	async function loadEvents(year: string) {
		loading = true;
		explorerError = null;
		selectedEvent = null;
		photos = [];
		selectedPhoto = null;
		previewUrl = null;
		try {
			const res = await fetch(`${data.apiURL}/api/photos/available/${encodeURIComponent(year)}`);
			if (!res.ok) throw new Error(`Failed to load events: ${res.status}`);
			events = await res.json();
			expandedEvents.add(year);
		} catch (e: any) {
			explorerError = e.message;
		} finally {
			loading = false;
		}
	}

	// Load photos in an event
	async function loadPhotos(year: string, event: string) {
		loading = true;
		explorerError = null;
		selectedPhoto = null;
		previewUrl = null;
		try {
			const res = await fetch(`${data.apiURL}/api/photos/available/${encodeURIComponent(year)}/${encodeURIComponent(event)}`);
			if (!res.ok) throw new Error(`Failed to load photos: ${res.status}`);
			photos = await res.json();
		} catch (e: any) {
			explorerError = e.message;
		} finally {
			loading = false;
		}
	}

	// Select a photo to preview
	function selectPhoto(filename: string) {
		selectedPhoto = filename;
		previewUrl = `${data.apiURL}/api/photos/available/${encodeURIComponent(selectedYear!)}/${encodeURIComponent(selectedEvent!)}/${encodeURIComponent(filename)}/preview`;
	}

	// Toggle year expand/collapse
	function toggleYear(year: string) {
		if (expandedYears.has(year)) {
			expandedYears.delete(year);
		} else {
			expandedYears.add(year);
		}
		expandedYears = new Set(expandedYears);
	}

	// Toggle event expand/collapse
	function toggleEvent(event: string) {
		if (expandedEvents.has(event)) {
			expandedEvents.delete(event);
		} else {
			expandedEvents.add(event);
		}
		expandedEvents = new Set(expandedEvents);
	}

	// Add photo from explorer
	async function handleAddPhoto() {
		if (!selectedPhoto || !addTitle.trim()) return;
		
		const filename = `${selectedYear}/${selectedEvent}/${selectedPhoto}`;
		const formData = new FormData();
		formData.append('filename', filename);
		formData.append('title', addTitle.trim());
		formData.append('sortOrder', String(addSortOrder));

		const res = await fetch(`${data.apiURL}/api/photos`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				filename,
				title: addTitle.trim(),
				sortOrder: addSortOrder
			})
		});

		if (!res.ok) {
			const body = await res.text();
			addFailure = body;
			addSuccess = null;
		} else {
			const result = await res.json();
			addSuccess = result.id;
			addFailure = null;
			// Reset form
			addTitle = '';
			addSortOrder = 0;
			selectedPhoto = null;
			previewUrl = null;
			// Reload photos list
			await loadPhotos(selectedYear!, selectedEvent!);
		}
	}

	// Regenerate all thumbnails
	async function handleRegenerate() {
		regenerating = true;
		regenerateFailure = null;
		regenerateSuccess = false;
		try {
			const res = await fetch(`${data.apiURL}/api/photos/regenerate-thumbnails`, {
				method: "POST"
			});
			if (!res.ok) {
				const body = await res.text();
				regenerateFailure = body;
			} else {
				regenerateSuccess = true;
			}
		} catch (e: any) {
			regenerateFailure = e.message;
		} finally {
			regenerating = false;
		}
	}

	// Initial load
	loadYears();
</script>

<section id="photos-editor" class="flex flex-col items-center gap-4 py-5 mt-16">
	<div class="flex w-full max-w-7xl items-center justify-between px-3">
		<h1 class="font-code text-lg">Manage Photo Gallery</h1>
		<div class="flex items-center gap-2">
			<button 
				class={buttonVariants({ size: 'sm', variant: 'outline' })} 
				onclick={handleRegenerate}
				disabled={regenerating}
			>
				<RefreshCw class={"mr-2 w-4 h-4" + (regenerating ? ' animate-spin' : '')} />
				{regenerating ? 'Regenerating...' : 'Regenerate Thumbnails'}
			</button>
			<button 
				class={buttonVariants({ size: 'sm' })} 
				onclick={() => showAddForm = !showAddForm}
			>
				<Plus class="mr-2 w-4 h-4" /> {showAddForm ? 'Cancel' : 'Add Photo'}
			</button>
		</div>
	</div>

	<!-- Add Photo Form -->
	{#if showAddForm && selectedYear && selectedEvent}
		<Card class="mx-3 w-2xl shadow-lg border-primary/20">
			<CardContent class="p-6">
				<h2 class="text-lg font-semibold mb-4">Add Photo</h2>
				<div class="flex flex-col gap-4">
					<div class="flex items-center gap-4">
						<div class="w-32 h-32 bg-muted rounded-lg flex items-center justify-center overflow-hidden flex-shrink-0">
							{#if previewUrl}
								<img src={previewUrl} alt="Preview" class="w-full h-full object-cover" />
							{:else}
								<ImageIcon class="w-12 h-12 opacity-30" />
							{/if}
						</div>
						<div class="flex-grow">
							<label class="flex flex-col gap-1">
								Title
								<input 
									bind:value={addTitle} 
									placeholder="Photo title"
									class="rounded-md border p-2 text-sm text-muted-foreground" 
									onkeydown={(e) => e.key === 'Enter' && handleAddPhoto()}
								/>
							</label>
							<label class="flex flex-col gap-1 mt-2">
								Sort Order
								<input 
									type="number" 
									bind:value={addSortOrder} 
									class="rounded-md border p-2 text-sm text-muted-foreground" 
								/>
							</label>
						</div>
					</div>
					<p class="text-sm text-muted-foreground">
						Selected: {selectedYear}/{selectedEvent}/{selectedPhoto || 'No photo selected'}
					</p>
					<div class="flex items-center gap-2">
						<button 
							class={buttonVariants({ size: 'sm' })}
							onclick={handleAddPhoto}
							disabled={!selectedPhoto || !addTitle.trim()}
						>
							Save Photo
						</button>
						{#if addSuccess !== null}
							<Check color={'green'} />
						{:else if addFailure}
							<X color={'red'} /> {addFailure}
						{/if}
					</div>
				</div>
			</CardContent>
		</Card>
	{/if}

	<!-- Main Content -->
	<div class="flex w-full max-w-7xl gap-4 px-3">
		<!-- File Explorer Panel -->
		<Card class="w-80 flex-shrink-0">
			<CardHeader class="pb-3">
				<CardTitle class="text-base">Photo Library</CardTitle>
				<CardDescription>Browse your photo collection</CardDescription>
			</CardHeader>
			<CardContent class="p-0">
				{#if loading}
					<div class="p-4 text-center text-sm text-muted-foreground">Loading...</div>
				{:else if explorerError}
					<div class="p-4">
						<p class="text-sm text-destructive">{explorerError}</p>
						<button class={buttonVariants({ size: 'sm', variant: 'outline', class: 'mt-2 w-full' })} onclick={loadYears}>
							Retry
						</button>
					</div>
				{:else}
					<div class="overflow-y-auto max-h-[60vh]">
						{#each years as year}
							<div>
								<button 
									class="w-full flex items-center gap-2 px-3 py-2 text-sm hover:bg-muted transition-colors"
									onclick={() => toggleYear(year)}
								>
									{#if expandedYears.has(year)}
										<ChevronDown class="w-4 h-4" />
									{:else}
										<ChevronRight class="w-4 h-4" />
									{/if}
									<Folder class="w-4 h-4" />
									<span>{year}</span>
								</button>
								{#if expandedYears.has(year)}
									<div class="ml-4 border-l border-muted pl-2">
										{#each events as event}
											<button 
												class="w-full flex items-center gap-2 px-3 py-2 text-sm hover:bg-muted transition-colors"
												onclick={() => {
													if (selectedYear !== year) {
														selectedYear = year;
														loadEvents(year);
													}
													if (selectedEvent !== event) {
														selectedEvent = event;
														loadPhotos(year, event);
													}
													toggleEvent(event);
												}}
											>
												{#if expandedEvents.has(event)}
													<ChevronDown class="w-4 h-4" />
												{:else}
													<ChevronRight class="w-4 h-4" />
												{/if}
												<Folder class="w-4 h-4" />
												<span class="truncate">{event}</span>
											</button>
											{#if expandedEvents.has(event) && selectedEvent === event}
												<div class="ml-4 border-l border-muted pl-2 max-h-64 overflow-y-auto">
													{#each photos as photo}
														<button 
															class="w-full flex items-center gap-2 px-3 py-1.5 text-sm hover:bg-muted transition-colors"
															onclick={() => selectPhoto(photo)}
														>
															<ImageIcon class="w-4 h-4" />
															<span class="truncate">{photo}</span>
														</button>
													{/each}
													{#if photos.length === 0}
														<div class="px-3 py-2 text-xs text-muted-foreground">No photos</div>
													{/if}
												</div>
											{/if}
										{/each}
										{#if events.length === 0}
											<div class="px-3 py-2 text-xs text-muted-foreground">No events</div>
										{/if}
									</div>
								{/if}
							</div>
						{/each}
						{#if years.length === 0}
							<div class="p-4 text-sm text-muted-foreground">No years found</div>
						{/if}
					</div>
				{/if}
			</CardContent>
		</Card>

		<!-- Gallery Management Panel -->
		<div class="flex-grow">
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-base font-semibold">Gallery Photos ({data.photos.length})</h2>
			</div>

			{#if form?.failure && !form?.id}
				<div class="mb-4 p-3 bg-destructive/10 text-destructive text-sm rounded-md">{form.message}</div>
			{/if}

			{#if regenerateFailure}
				<div class="mb-4 p-3 bg-destructive/10 text-destructive text-sm rounded-md">{regenerateFailure}</div>
			{:else if regenerateSuccess}
				<div class="mb-4 p-3 bg-green-500/10 text-green-500 text-sm rounded-md">Thumbnails regenerated successfully</div>
			{/if}

			<div class="flex flex-col gap-4">
				{#each data.photos as photo}
					<Card class="shadow-md overflow-hidden">
						<CardContent class="p-0 flex flex-col sm:flex-row">
							<!-- Photo Preview -->
							<div class="w-full sm:w-48 h-48 bg-muted flex items-center justify-center overflow-hidden relative group">
								{#if photo.mediumPath}
									<img 
										src={`${data.apiURL}/thumbnails/${photo.mediumPath}`} 
										alt={photo.altText || photo.title}
										class="w-full h-full object-cover"
									/>
								{:else if photo.filePath}
									<img 
										src={`${data.apiURL}/api/photos/${photo.id}/image`} 
										alt={photo.altText || photo.title}
										class="w-full h-full object-cover"
									/>
								{:else}
									<ImageIcon class="w-12 h-12 opacity-30" />
								{/if}
								{#if !photo.visible}
									<div class="absolute inset-0 bg-background/50 flex items-center justify-center">
										<span class="text-xs font-bold opacity-50">HIDDEN</span>
									</div>
								{/if}
							</div>
							
							<!-- Photo Info & Edit Form -->
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
											<span class="text-xs font-bold opacity-50">SORT</span>
											<input type="number" name="sortOrder" value={photo.sortOrder} class="rounded border p-1 text-sm text-muted-foreground" />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">CAMERA</span>
											<input value={photo.camera || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">LENS</span>
											<input value={photo.lens || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">APERTURE</span>
											<input value={photo.aperture || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">SHUTTER</span>
											<input value={photo.shutterSpeed || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">ISO</span>
											<input value={photo.iso || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">DATE</span>
											<input value={photo.dateTaken || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
										<label class="flex flex-col gap-1">
											<span class="text-xs font-bold opacity-50">LOCATION</span>
											<input value={photo.location || ''} class="rounded border p-1 text-sm text-muted-foreground" readonly />
										</label>
									</div>

									<div class="flex items-center justify-between mt-2">
										<div class="flex items-center gap-2">
											<input type="checkbox" name="visible" checked={photo.visible} />
											<span class="text-xs text-muted-foreground">Visible</span>
										</div>
										
										<div class="flex items-center gap-2">
											<span class="text-xs text-muted-foreground truncate max-w-[200px]" title={photo.sourcePath}>
												{photo.sourcePath}
											</span>
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
							
							<!-- Delete Button -->
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
		</div>
	</div>
</section>

<style>
	.animate-spin {
		animation: spin 1s linear infinite;
	}
	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}
</style>
