<script lang="ts">
	interface SearchResult {
		title: string;
		path: string /*fix how interface should look to correspond with backend*/;
		type: string;
		source: string;
	}

	let query: string = '';
	let submittedQuery = '';
	let results: SearchResult[] = [];
	let searched: boolean = false;

	async function search(): Promise<void> {
		searched = true;
		submittedQuery = query;
		/*
        const res = await fetch(`http://localhost:8080/search?q=${query}`); /*change to url we are using for backend
        results = await res.json() as SearchResult[];
        */
		results = [
			{
				title: 'result 1',
				path: 'home/Path2D',
				type: 'File',
				source: 'OSPP'
			},
			{
				title: 'result 2',
				path: 'http://ddsd',
				type: 'Webbsite',
				source: 'Wikipedia'
			}
		];
	}
</script>

<main style="max-width: 800px; margin: 2rem auto; padding: 1rem;">

	<div style="display: flex; gap: 0.5rem; margin-bottom: 2rem;">
		<input 
			bind:value={query} 
			type="text" 
			placeholder="Write your search here!" 
			style="flex: 1; padding: 0.75rem; font-size: 1rem; background-color: #F0EEEE; border-radius: 6px; border: none;"
		/>

		<button on:click={search} style="padding: 0.75rem 1rem; font-size: 1rem; background-color: #AEC6DF; color: black; border: none; border-radius: 6px;"> 
			Search 
		</button>
	</div>

	{#if searched == true && results.length > 0}
			{#each results as res}
				<div style="background: white; border: 1px solid #ddd; border-radius: 8px; padding: 1rem; margin-bottom: 1rem;">
					<div style="display: flex; justify-content: space-between; align-items: center;">
						<h2 style="font-size: 1.2rem;">{res.title}</h2>
						<small>
							{res.type}: {res.source} 
						</small>
					</div>
					<!-- Add <p> with description here is we decide to show it-->
				</div>
			{/each}
	{:else if searched == true}
		<p>no results for: {submittedQuery}</p>
	{/if}
</main>