<script lang="ts">
	import { onMount } from 'svelte';
	import { showFiles, showWebpages, showAllResults, maxResults } from '$lib/stores/settings';
	import { get } from 'svelte/store';

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

	let searchInput: HTMLInputElement;

	async function search(): Promise<void> {
		if (query.length > 0) {
			submittedQuery = query;
			const res = await fetch(`http://localhost:8080/search?q=${query}`);
			const json = (await res.json()) as SearchResponse;
			let filteredResults = json.Results;

			filteredResults = filteredResults.filter(
				(res) => (res.Source === 1 && get(showWebpages)) || (res.Source !== 1 && get(showFiles))
			);

			if (!get(showAllResults)) {
				filteredResults = filteredResults.slice(0, get(maxResults));
			}
			console.log(results);
			results = filteredResults;
			searched = true;
		} else {
			searched = false;
		}
	}

	async function refreshSearch(): Promise<void> {
		if (submittedQuery.length > 0) {
			query = submittedQuery;
			await search();
		}
	}

	async function downloadFile(path: string): Promise<void> {
		fetch(`http://localhost:8080/download?q=${path}`, {
			method: 'GET'
		})
			.then((response) => {
				return response.blob();
			})
			.then((blob) => {
				const urlObject = window.URL.createObjectURL(blob);
				const a = document.createElement('a');

				a.href = urlObject;
				document.body.appendChild(a);
				a.click();
				document.body.removeChild(a);

				window.URL.revokeObjectURL(urlObject);
			});
	}

	onMount(() => {
		searchInput.focus();
	});
</script>

<main style="max-width: 800px;">
	<div id="searchDiv">
		<input
			bind:value={query}
			bind:this={searchInput}
			type="text"
			placeholder="Write your search here!"
			on:keyup={search}
		/>

		<button on:click={search} id="searchButton"> Search </button>

		<button on:click={refreshSearch} class="round-button" title="Refresh results"> ↻ </button>
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
						<div id="resultDiv">
							<h3>{res.Path}</h3>
						</div>
						<div id="resultInfo">
							<p class="searchInfo">
								Website: {res.Path}
							</p>
							<p class="searchInfo">
								Relevance: {res.Score.toFixed(4)}
							</p>
						</div>
					</div>
				</a>
			{:else}
				<div id="resultBox">
					<div id="resultDiv">
						<h3 style="margin:0">
							{res.Path.replace(/^.*[\\\/]/, '')}
						</h3>
						<button on:click={() => downloadFile(res.Path)} id="downloadButton"> Download </button>
					</div>
					<div id="resultInfo">
						<p class="searchInfo">
							Local file path: {res.Path}
						</p>
						<p class="searchInfo">
							Relevance: {res.Score.toFixed(4)}
						</p>
					</div>
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
		align-items: center;
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
		border: 2px solid #f0eeee;
		border-radius: 8px;
		padding: 1rem;
		margin-bottom: 1rem;
	}

	#resultDiv {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-top: 0.8rem;
	}

	#resultInfo {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-top: 1.5rem;
		margin-bottom: 1rem;
	}

	.searchInfo {
		font-size: rem;
		margin: 0;
	}

	.round-button {
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 50%;
		background-color: #f0eeee;
		cursor: pointer;
		font-size: 1.2rem;
		line-height: 1;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.round-button:hover {
		background-color: #e2e2e2;
	}
</style>

