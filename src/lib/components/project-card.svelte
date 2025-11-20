<script lang="ts">
	import { Card, CardContent, CardFooter, CardTitle } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Puzzle, Globe } from '@lucide/svelte';

	export let title;
	export let deploymentLink = '';
	export let technologies;
	export let description;
	export let image = '';
	export let altText = 'Project image';

	const iconMap = {
		Puzzle,
		Globe
	};

	type Icon = keyof typeof iconMap;

	export let icon: Icon = 'Puzzle';

	$: IconComponent = iconMap[icon] || Puzzle;

	// Function to bold technologies found in the description
	$: highlightedDescription = highlightTechnologies(description, technologies);
	function highlightTechnologies(text: string, techs: string[]) {
		if (!text || !techs?.length) return text;

		// Sort longer tech names first to avoid partial matches (e.g. "C" in "C++")
		const sortedTechs = [...techs].sort((a, b) => b.length - a.length);

		// Create a regex to match any technology word (case-insensitive)
		const regex = new RegExp(`\\b(${sortedTechs.join('|')})\\b`, 'gi');

		return text.replace(regex, (match) => `<span class="font-semibold">${match}</span>`);
	}
</script>

<Card class="mx-3 max-w-3xl shadow-lg">
	<CardTitle
		class="font-code flex min-h-8 flex-row flex-wrap items-center justify-between gap-2 px-6 text-lg sm:text-xl"
	>
		<h1 class="flex items-center gap-2 font-normal text-wrap">
			<svelte:component this={IconComponent} size={22} class="mb-0.5" />
			{title}
		</h1>
		<!-- show a 'Try it!' button when a deployment link is passed -->
		{#if deploymentLink != ''}
			<Button size="sm" class="font-code">
				<a href={deploymentLink} target="_blank" rel="noopener noreferrer"> Try it! </a>
			</Button>
		{/if}
	</CardTitle>
	<CardContent class={image != '' ? 'grid gap-6 sm:grid-cols-2' : 'flex'}>
		{#if image}
			<img
				src={image}
				alt={altText}
				class="aspect-square max-w-full rounded-lg object-cover sm:max-w-full"
			/>
		{/if}
		<p>{@html highlightedDescription}</p>
	</CardContent>
	<CardFooter class="flex flex-wrap gap-2 border-t">
		{#each technologies as tech}
			<Badge variant="secondary">{tech}</Badge>
		{/each}
	</CardFooter>
</Card>
