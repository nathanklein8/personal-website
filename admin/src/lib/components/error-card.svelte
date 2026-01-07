<script lang="ts">
	import { Alert, AlertDescription, AlertTitle } from '$lib/components/ui/alert';
	import { CloudAlert } from '@lucide/svelte';

	export let title = 'Uh Oh!';
	export let description;
</script>

<Alert
	variant="destructive"
	class="w-fit max-w-[85vw] md:max-w-[70vw] lg:max-w-[55vw] xl:max-w-[40vw]"
>
	<CloudAlert />
	<AlertTitle class="font-code text-lg">{title}</AlertTitle>
	<AlertDescription>
		<p>{description}</p>
	</AlertDescription>
</Alert>
