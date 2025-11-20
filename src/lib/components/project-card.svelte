<script lang="ts">
	import { Card, CardContent, CardFooter, CardTitle } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Puzzle, Globe } from '@lucide/svelte';

	const iconMap = {
		Puzzle,
		Globe
	};

	type Icon = keyof typeof iconMap;

	interface Props {
		icon: Icon;
		title: string;
		description: string;
		technologies: string[];
		deploymentLink?: string;
		image?: string;
		altText?: string;
	}

	let {
		icon,
		title,
		description,
		technologies,
		deploymentLink = '',
		image = '',
		altText = 'Project image'
	}: Props = $props();

	const IconComponent = iconMap[icon];

	// Function to bold technologies found in the description
	function highlightTechnologies(text: string, techs: string[]) {
		
		// Create a regex to match any technology word (case-insensitive)
		const regex = new RegExp(`\\b(${techs.join('|')})\\b`, 'gi');

		return text.replace(
			regex,
			(match) => `<span class="font-semibold text-foreground/80">${match}</span>`
		);
	}
	// const highlightedDescription = highlightTechnologies(description, technologies);
</script>

<Card class="mx-3 max-w-3xl shadow-lg">
	<CardTitle
		class="font-code flex min-h-8 flex-row flex-wrap items-center justify-between gap-2 px-6 text-lg sm:text-xl"
	>
		<h1 class="flex items-center gap-2 font-normal text-wrap">
			<IconComponent class="mb-0.5"/>
			{title}
		</h1>
		<!-- show a 'Try it!' button when a deployment link is passed -->
		{#if deploymentLink != ''}
			<Button size="sm" class="font-code">
				<a href={deploymentLink} target="_blank" rel="noopener noreferrer"> Try it! </a>
			</Button>
		{/if}
	</CardTitle>
	<CardContent class={image != '' ? 'grid gap-6 sm:grid-cols-2' : ''}>
		{#if image}
			<img
				src={image}
				alt={altText}
				class="aspect-square max-w-full rounded-lg object-cover sm:max-w-full"
			/>
		{/if}
		<p class="text-muted-foreground">
			{@html highlightTechnologies(description, technologies)}
		</p>
	</CardContent>
	<CardFooter class="flex flex-wrap gap-2 border-t">
		{#each technologies as tech}
			<Badge variant="secondary">{tech}</Badge>
		{/each}
	</CardFooter>
</Card>
