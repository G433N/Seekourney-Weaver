<script lang="ts">

	interface SearchResult {
		Path: string;
		Score: int;
		Source: int;
	} /*fix how interface should look to correspond with backend*/

	interface SearchResponse {
		Query: string;
		Results: SearchResult[];
	}


	let query: string = '';
	let submittedQuery = '';
	let results: SearchResult[] = [];
	let searched: boolean = false;

	async function search(): Promise<void> {
		searched = true;
		submittedQuery = query;
        const res = await fetch(`http://localhost:8080/search?q=${query}`);
        results = await res.json() as SearchResult[];
		console.log(results);
		results = results.Results;
		// results = [
		// 	{
		// 		title: 'result 1',
		// 		path: 'https://en.wikipedia.org/wiki/Bear',
		// 		type: 'File',
		// 		source: 'OSPP',
		// 		desc: 'Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden'
		// 	},
		// 	{
		// 		title: 'result 2',
		// 		path: 'https://en.wikipedia.org/wiki/Bear',
		// 		type: 'Webbsite',
		// 		source: 'Wikipedia',
		// 		desc: 'Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden'
		// 	}
		// ];
	}
</script>

<main style="max-width: 800px;">

	<div id="searchDiv">
		<input 
			bind:value={query} 
			type="text" 
			placeholder="Write your search here!"
		/>

		<button on:click={search} id="searchButton"> 
			Search 
		</button>
	</div>

	{#if searched == true && results.length > 0}
			{#each results as res}
				<a 
					href={res.path} 
					target="_blank" 
					rel="noopener noreferrer"
					style="display: block; text-decoration: none; color: inherit;"
				>
					<div id="resultBox">
						<div  id="resultDiv">
							<h3 style="font-size: 1.4rem;">{res.Path}</h3>
							<small>
								{res.Score}, {res.Source} 
							</small>
						</div>
						<p style="color: #4E4E4E;">{res.desc}</p> 
					</div>
				</a>
			{/each}
	{:else if searched == true}
		<p style="font-size: 1.2rem;">no results for: {submittedQuery}</p>
	{/if}
</main>

<style>
	#searchDiv {
		display: flex; 
		gap: 0.5rem; 
		margin-bottom: 2rem;
	}

	#searchButton {
		padding: 0.75rem 1rem; 
		font-size: 1rem; 
		background-color: #AEC6DF; 
		color: black; 
		border: none; 
		border-radius: 6px;
	}

	#resultBox {
		background: white; 
		border: 2px solid #F0EEEE; 
		border-radius: 8px; 
		padding: 1rem; 
		margin-bottom: 1rem;
	}

	#resultDiv {
		display: flex; 
		justify-content: space-between; 
		align-items: center;
	}
</style>
