<script lang="ts">
	import { writable } from 'svelte/store';

	const API_URL = import.meta.env.VITE_API_URL;

	let newNote = '';
	let incrementId: number | null = null;

	const createResponse = writable<any>(null);
	const incrementResponse = writable<any>(null);

	// Create a new test item
	async function handleCreate() {
		console.log('creating' + newNote);
		console.log('calling ' + `${API_URL}/api/test`);
		const res = await fetch(`${API_URL}/api/test`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ note: newNote })
		});

		if (res.ok) {
			const item = await res.json();
			createResponse.set(item);
			newNote = '';
		} else {
			const text = await res.text();
			createResponse.set({ error: text });
		}
	}

	// Increment count for a given test item
	async function handleIncrement() {
		console.log('incrementing' + incrementId);
		console.log('calling ' + `${API_URL}/api/test/${incrementId}/increment`);
		if (incrementId === null) return;

		const res = await fetch(`${API_URL}/api/test/${incrementId}/increment`, {
			method: 'POST'
		});

		if (res.ok) {
			const updated = await res.json();
			incrementResponse.set(updated);
		} else {
			const text = await res.text();
			incrementResponse.set({ error: text });
		}
	}
</script>

<section class="flex h-screen flex-col items-center justify-center bg-gray-50 gap-6 p-4">

	<div>API URL: {API_URL}</div>

	<div class="flex flex-col items-center gap-2">
		<h2 class="text-xl font-semibold">Create Test Item</h2>
		<input
			type="text"
			bind:value={newNote}
			placeholder="Enter note"
			class="border px-2 py-1 rounded"
		/>
		<button on:click={handleCreate} class="bg-blue-500 text-white px-4 py-1 rounded">
			Create
		</button>
		{#if $createResponse}
			<pre class="bg-gray-100 p-2 rounded w-80 overflow-x-auto">{JSON.stringify($createResponse, null, 2)}</pre>
		{/if}
	</div>

	<div class="flex flex-col items-center gap-2">
		<h2 class="text-xl font-semibold">Increment Count</h2>
		<input
			type="number"
			bind:value={incrementId}
			placeholder="Enter ID"
			class="border px-2 py-1 rounded"
		/>
		<button on:click={handleIncrement} class="bg-green-500 text-white px-4 py-1 rounded">
			Increment
		</button>
		{#if $incrementResponse}
			<pre class="bg-gray-100 p-2 rounded w-80 overflow-x-auto">{JSON.stringify($incrementResponse, null, 2)}</pre>
		{/if}
	</div>

</section>
