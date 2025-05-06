<script lang="ts">

	interface SearchResult {
		Path: string;
		Score: number;
		Source: number;
	}

	interface SearchResponse {
		Query: string;
		Results: SearchResult[];
	}

	let query: string = '';
	let submittedQuery: string = '';
	let results: SearchResult[] = [];
	let searched: boolean = false;

	async function search(): Promise<void> {
		if (query.length > 0)
		{
			searched = true;
			submittedQuery = query;
			const res = await fetch(`http://localhost:8080/search?q=${query}`);
			const json = await res.json() as SearchResponse;
			results = json.Results;
			console.log(results);
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
		else 
		{
			searched = false;
		}
	}

	// TODO: test if it works with branch search-and-download
	async function downloadFile(path: string): Promise<void> {
		fetch(`http://localhost:8080/download?q=${path}`, {
		method: "GET"
		})
		.then(response => {
		return response.blob();
		})
		.then(blob => {
		const urlObject = window.URL.createObjectURL(blob);
		const a = document.createElement("a");

		a.href = urlObject;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);

		window.URL.revokeObjectURL(urlObject);
		});
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
				{#if res.Source == 1}
				<a 
					href={res.Path} 
					target="_blank" 
					rel="noopener noreferrer"
					style="display: block; text-decoration: none; color: inherit;"
				>
					<div id="resultBox">
						<div  id="resultDiv">
							<h3>{res.Path}</h3>
						</div>
						<div id=resultInfo>
							<small> 
								Website: {res.Path}
							</small>
							<small>
								Relevance: {res.Score.toFixed(4)}
							</small>
						</div>
						<!-- <p style="color: #4E4E4E;">{res.desc}</p> -->
					</div>
				</a>
				{:else}
					<div id="resultBox">
						<div id="resultDiv">
							<h3>{res.Path.replace(/^.*[\\\/]/, '')}</h3>
							<button on:click={() => downloadFile(res.Path)} id="downloadButton">
								Download
							</button>
						</div>
						<div id=resultInfo>
							<small> 
								Local file path: {res.Path} 
							</small>
							<small> 
								Relevance: {res.Score.toFixed(4)} 
							</small>
						</div>
						<!-- <p style="color: #4E4E4E;">{res.desc}</p> -->
					</div>
				{/if}
				
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
		font-weight: 500;
	}

	#downloadButton {
		padding: 0.2rem 0.6rem; 
		font-size: 1rem;
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

	#resultInfo {
		display: flex; 
		justify-content: space-between; 
		align-items: center;
	}
</style>
