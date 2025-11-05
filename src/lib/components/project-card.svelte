<script lang="ts">
	import { Card, CardContent, CardFooter, CardTitle } from './ui/card';
	import { Button } from './ui/button';
	import { Badge } from './ui/badge';
    import { Puzzle, Globe } from '@lucide/svelte';

	export let title;
	export let deploymentLink = '';
	export let technologies;
    export let description;
    export let image = '';
    export let altText = 'Project image';
    
    const iconMap = {
		Puzzle,
		Globe,
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

		return text.replace(regex, (match) => `<strong>${match}</strong>`);
	}
</script>

<Card
	class="max-w-[95vw] shadow-lg md:max-w-[80vw] lg:max-w-[65vw] xl:max-w-[50vw]"
>
	<CardTitle class="font-code flex grid flex-row items-center gap-3 px-6 text-xl sm:grid-cols-2">
		<h1 class="flex items-center gap-2">
			<svelte:component this={IconComponent} size={22} />
			{title}
		</h1>
		{#if deploymentLink != ''}
			<div class="flex">
                <span class="flex md:grow"></span>
				<Button variant="green">
					<a
						class="font-code text-md md:text-lg"
						href={deploymentLink}
						target="_blank"
						rel="noopener noreferrer"
					>
						check it out
					</a>
				</Button>
			</div>
		{/if}
	</CardTitle>
	<CardContent class={image != "" ? "grid gap-6 sm:grid-cols-2" : "flex"}>
        {#if image}
        <img
			src={image}
			alt={altText}
			class="max-w-full sm:max-w-full rounded-lg object-cover aspect-square"
		/>
        {/if}
		<p class="text-md text-muted-foreground">
			{@html highlightedDescription}
		</p>
	</CardContent>
	<CardFooter class="flex flex-wrap gap-2 border-t">
		{#each technologies as tech}
			<Badge variant="secondary">{tech}</Badge>
		{/each}
	</CardFooter>
</Card>
