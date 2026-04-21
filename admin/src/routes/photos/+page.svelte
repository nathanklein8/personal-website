<script lang="ts">
	import type { ActionData, PageData } from './$types';
	import { getAllPhotos } from '@nk/shared/server/backend';

	let { data, form, errors }: PageData & { form: ActionData; errors: ActionData } = $props();

	// State management
	let showBrowser = $state(false);
	let selectedYear = $state<string | null>(null);
	let selectedEvent = $state<string | null>(null);
	let selectedPhoto = $state<{ year: string; event: string; filename: string } | null>(null);
	let isLoading = $state(false);
	let addLoading = $state(false);
	let addSuccess = $state<string | null>(null);
	let addError = $state<string | null>(null);

	// Data from API
	let years = $state<string[]>([]);
	let events = $state<string[]>([]);
	let availablePhotos = $state<string[]>([]);
	let addedPhotos = $state<Photo[]>([]);

	// Caption input
	let captionInput = $state('');

	// Load initial data
	async function loadData() {
		isLoading = true;
		try {
			const [yearsRes, photosRes] = await Promise.all([
				fetch('/api/photos/available'),
				getAllPhotos()
			]);

			if (yearsRes.ok) {
				years = await yearsRes.json();
			}

			addedPhotos = Array.isArray(photosRes) ? photosRes : [];
		} catch (e) {
			console.error('Failed to load initial data:', e);
		} finally {
			isLoading = false;
		}
	}

	await loadData();

	// Year selection
	async function selectYear(year: string) {
		selectedYear = year;
		selectedEvent = null;
		availablePhotos = [];

		const res = await fetch(`/api/photos/available/${encodeURIComponent(year)}`);
		if (res.ok) {
			events = await res.json();
		}
	}

	// Event selection
	async function selectEvent(event: string) {
		selectedEvent = event;

		const res = await fetch(
			`/api/photos/available/${encodeURIComponent(selectedYear!)}/${encodeURIComponent(event)}`
		);
		if (res.ok) {
			availablePhotos = await res.json();
		}
	}

	// Photo selection (opens popover)
	function selectPhoto(filename: string) {
		const nameWithoutExt = filename.replace(/\.[^/.]+$/, '');
		captionInput = nameWithoutExt;
		selectedPhoto = {
			year: selectedYear!,
			event: selectedEvent!,
			filename
		};
		addSuccess = null;
		addError = null;
	}

	// Close popover
	function closePopover() {
		selectedPhoto = null;
		captionInput = '';
		addSuccess = null;
		addError = null;
	}

	// Add photo
	async function addPhoto() {
		if (!selectedPhoto) return;

		addLoading = true;
		addError = null;
		addSuccess = null;

		try {
			// Calculate sort order: current count + 1
			const sortOrder = addedPhotos.length + 1;

			const res = await fetch('/api/photos', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					filename: `${selectedPhoto.year}/${selectedPhoto.event}/${selectedPhoto.filename}`,
					title: captionInput,
					sortOrder
				})
			});

			if (!res.ok) {
				const body = await res.text();
				addError = body || `Backend error: ${res.status}`;
				return;
			}

			const result = await res.json();
			addSuccess = `Photo added successfully (ID: ${result.id})`;

			// Refresh added photos list
			addedPhotos = await getAllPhotos();

			// Close popover after a brief delay
			setTimeout(closePopover, 1500);
		} catch (e) {
			addError = e instanceof Error ? e.message : 'Failed to add photo';
		} finally {
			addLoading = false;
		}
	}

	// Navigate back
	function goBackToYears() {
		selectedYear = null;
		selectedEvent = null;
		availablePhotos = [];
		events = [];
	}

	function goBackToEvents() {
		selectedEvent = null;
		availablePhotos = [];
	}

	// Get preview URL for a photo
	function getPreviewUrl(year: string, event: string, filename: string): string {
		return `/api/photos/available/${encodeURIComponent(year)}/${encodeURIComponent(event)}/${encodeURIComponent(filename)}/preview`;
	}

	// Get the display name for an event (strip the "N: " prefix)
	function getEventDisplay(name: string): string {
		return name.replace(/^[0-9]+:\s*/, '');
	}

	// Get the filename without path and extension
	function getPhotoDisplay(filename: string): string {
		const parts = filename.split('/');
		const base = parts[parts.length - 1];
		return base.replace(/\.[^/.]+$/, '');
	}

	// Get photo sort order for the next photo
	function getNextSortOrder(): number {
		return addedPhotos.length + 1;
	}
</script>

<svelte:head>
	<title>Photo Editor - Personal Website</title>
</svelte:head>

<section class="flex flex-col items-center gap-6 py-8 px-4 max-w-6xl mx-auto">
	<!-- Header -->
	<div class="flex items-center justify-between w-full">
		<h1 class="text-2xl font-bold">Photo Manager</h1>
		<button
			class="btn btn-primary"
			onclick={() => {
				showBrowser = !showBrowser;
				if (!showBrowser) {
					// Reset browser state when closing
					selectedYear = null;
					selectedEvent = null;
					availablePhotos = [];
					events = [];
				}
			}}
		>
			{showBrowser ? 'Close Browser' : 'Add Photo'}
		</button>
	</div>

	<!-- File Browser Section -->
	{#if showBrowser}
		<div class="w-full space-y-4">
			<!-- Breadcrumb -->
			{#if selectedYear || selectedEvent}
				<nav class="flex items-center gap-2 text-sm text-muted-foreground">
					<button
						class="hover:text-foreground transition-colors"
						onclick={goBackToYears}
					>
						Years
					</button>
					{#if selectedEvent}
						<span>/</span>
						<button
							class="hover:text-foreground transition-colors"
							onclick={goBackToEvents}
						>
							{getEventDisplay(selectedEvent!)}
						</button>
						<span>/</span>
						<span class="text-foreground">Photos</span>
					{:else}
						<span>/</span>
						<span class="text-foreground">Events</span>
					{/if}
				</nav>
			{/if}

			<!-- Loading State -->
			{#if isLoading}
				<div class="flex items-center justify-center py-12">
					<p class="text-muted-foreground">Loading...</p>
				</div>
			{/if}

			<!-- Year List -->
			{#if !selectedYear && years.length > 0}
				<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
					{#each years.sort((a, b) => b - a) as year}
						<button
							class="card card-bordered p-4 text-center hover:shadow-md transition-all cursor-pointer"
							onclick={() => selectYear(year)}
						>
							<span class="text-lg font-semibold">{year}</span>
						</button>
					{/each}
				</div>

			<!-- Event List -->
			{:else if selectedYear && !selectedEvent && events.length > 0}
				<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
					{#each events as event}
						<button
							class="card card-bordered p-4 text-left hover:shadow-md transition-all cursor-pointer"
							onclick={() => selectEvent(event)}
						>
							<span class="font-medium">{getEventDisplay(event)}</span>
						</button>
					{/each}
				</div>

			<!-- Photo Grid -->
			{:else if selectedEvent && availablePhotos.length > 0}
				<div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-3">
					{#each availablePhotos as filename}
						<button
							class="group relative aspect-square rounded-lg overflow-hidden border border-border hover:border-primary transition-all cursor-pointer"
							onclick={() => selectPhoto(filename)}
						>
							<img
								src={getPreviewUrl(selectedYear!, selectedEvent!, filename)}
								alt={getPhotoDisplay(filename)}
								class="w-full h-full object-cover"
								loading="lazy"
							/>
							<div class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-all flex items-end">
								<span class="text-xs text-white p-1.5 opacity-0 group-hover:opacity-100 transition-opacity truncate w-full">
									{getPhotoDisplay(filename)}
								</span>
							</div>
						</button>
					{/each}
				</div>

			<!-- Empty State -->
			{:else if !isLoading && availablePhotos.length === 0 && selectedEvent}
				<div class="flex items-center justify-center py-12">
					<p class="text-muted-foreground">No photos found in this event.</p>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Photo Popover (Dialog) -->
	{#if selectedPhoto}
		<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onclick={closePopover}>
			<div
				class="bg-background border border-border rounded-lg shadow-lg p-6 w-full max-w-md mx-4 space-y-4"
				onclick={(e) => e.stopPropagation()}
			>
				<!-- Preview Image -->
				<div class="relative aspect-video rounded-lg overflow-hidden bg-muted">
					<img
						src={getPreviewUrl(selectedPhoto.year, selectedPhoto.event, selectedPhoto.filename)}
						alt={captionInput}
						class="w-full h-full object-contain"
					/>
				</div>

				<!-- Caption Input -->
				<div class="space-y-2">
					<label for="caption" class="text-sm font-medium">Caption</label>
					<input
						id="caption"
						type="text"
						bind:value={captionInput}
						class="input input-bordered w-full"
						placeholder="Enter caption..."
						onkeydown={(e) => {
							if (e.key === 'Enter') addPhoto();
						}}
					/>
				</div>

				<!-- Status Messages -->
				{#if addSuccess}
					<p class="text-sm text-green-600 dark:text-green-400">{addSuccess}</p>
				{/if}
				{#if addError}
					<p class="text-sm text-red-600 dark:text-red-400">{addError}</p>
				{/if}

				<!-- Actions -->
				<div class="flex justify-end gap-2 pt-2">
					<button
						class="btn btn-outline"
						onclick={closePopover}
						disabled={addLoading}
					>
						Close
					</button>
					<button
						class="btn btn-primary"
						onclick={addPhoto}
						disabled={addLoading || !captionInput.trim()}
					>
						{#if addLoading}
							<span class="loading loading-spinner loading-sm"></span>
							Adding...
						{:else}
							Add Photo
						{/if}
					</button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Added Photos List -->
	{#if addedPhotos.length > 0}
		<div class="w-full space-y-4 pt-6 border-t border-border">
			<h2 class="text-xl font-bold">Added Photos</h2>
			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
				{#each addedPhotos as photo}
					<div class="card card-bordered">
						<div class="card-body p-3">
							{#if photo.thumbnailPath}
								<img
									src={`/api/photos/available/${photo.sourcePath}/preview`}
									alt={photo.title}
									class="w-full h-32 object-cover rounded-md mb-2"
								/>
							{/if}
							<h3 class="text-sm font-medium truncate">{photo.title}</h3>
							{#if photo.camera}
								<p class="text-xs text-muted-foreground">{photo.camera}</p>
							{/if}
							<div class="flex gap-1 mt-1 flex-wrap">
								{#if photo.aperture}
									<span class="badge badge-sm badge-outline">f/{photo.aperture}</span>
								{/if}
								{#if photo.shutterSpeed}
									<span class="badge badge-sm badge-outline">{photo.shutterSpeed}</span>
								{/if}
								{#if photo.iso}
									<span class="badge badge-sm badge-outline">ISO {photo.iso}</span>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</section>

<style>
	.btn {
		@apply inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50;
		@apply h-10 px-4 py-2;
	}

	.btn-primary {
		@apply bg-primary text-primary-foreground hover:bg-primary/90;
	}

	.btn-outline {
		@apply border border-input bg-background hover:bg-accent hover:text-accent-foreground;
	}

	.card {
		@apply rounded-lg border bg-card text-card-foreground shadow-sm;
	}

	.card-bordered {
		@apply border border-border;
	}

	.card-body {
		@apply p-6;
	}

	.input {
		@apply flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50;
	}

	.input-bordered {
		@apply border border-border;
	}

	.badge {
		@apply inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2;
	}

	.badge-outline {
		@apply border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80;
	}

	.loading {
		@apply animate-spin;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
	}
</style>
