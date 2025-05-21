<script lang="ts">
	import {
		showFiles,
		showWebpages,
		showAllResults,
		maxResults,
		cpuDefault,
		maxCores,
		cpuCores,
		indexerList
	} from 'src/lib/stores/settings';

	import { get } from 'svelte/store';

	interface IndexerResult {
		Name: string;
		Id: number;
		Port: number;
	}

	let indexerInput: '';
	// let submittedIndexer: string = '';

	async function addIndexer(): Promise<void> {
		if (indexerInput.length > 0) {
			// submittedIndexer = indexerInput;
			//const res = await fetch(`http://localhost:8080/addIndexer?q=${submittedIndexer}`); //TODO: what name??
			//const json = await res.json() as IndexerResult;
			//indexerList = [...indexerList, json];

			let mockResult: IndexerResult = {
				Name: 'indexer',
				Id: Math.round(Math.random() * 100),
				Port: 1
			};

			indexerList = [...indexerList, mockResult];
			indexerInput = '';
			// TODO: unsure of how response should look, might not work
			// TODO: add some fix for duplicates?
		}
	}

	async function deleteIndexer(indexer: IndexerResult): Promise<void> {
		fetch(`http://localhost:8080/addIndexer?q=${indexer}`); //TODO: what name??

		indexerList = indexerList.filter((elem) => elem.Id !== indexer.Id);

		// TODO: unsure if we get a response?
	}
</script>

<main style="max-width: 600px;">
	<h2>Filter</h2>
	<div class="box">
		<h3>Show:</h3>
		<label class="toggle">
			<input type="checkbox" bind:checked={showFiles} />
			Files
		</label>
		<label class="toggle">
			<input type="checkbox" bind:checked={showWebpages} />
			Webpages
		</label>
	</div>

	<div class="box column">
		<label class="toggle">
			<input type="checkbox" bind:checked={showAllResults} />
			Show all results
		</label>

		{#if !showAllResults}
			<label class="max-label">
				<div class="inputDiv">
					Max results shown:
					<input type="number" bind:value={maxResults} />
				</div>
			</label>
		{/if}
	</div>

	<h2>Searching</h2>
	<div class="box column">
		<label class="toggle">
			<input type="checkbox" bind:checked={cpuDefault} />
			default CPU usage
		</label>
		{#if !cpuDefault}
			<div class="slider">
				<label for="cpuSlider">CPU cores used:</label>
				<input id="cpuSlider" type="range" min="1" max={maxCores} bind:value={cpuCores} />
				<span>{cpuCores}</span>
			</div>
		{/if}
	</div>

	<h2>Indexer</h2>
	<div class="box column">
		<label class="max-label">
			<div class="">
				Indexer path:
				<input type="text" bind:value={indexerInput} />
				<button id="indexerButton" onclick={() => addIndexer()}> Add </button>
			</div>
		</label>

		{#if indexerList.length > 0}
			{#each indexerList as indexer}
				<div class="inputDiv">
					<h3>
						{indexer.Name}, ID: {indexer.Id}, Port: {indexer.Port}
					</h3>
					<button id="deleteButton" onclick={() => deleteIndexer(indexer)}> Delete </button>
				</div>
			{/each}
		{/if}
	</div>
</main>

<style>
	.box {
		background: white;
		border: 2px solid #f0eeee;
		border-radius: 12px;
		padding: 1rem 1.5rem;
		margin-bottom: 1.5rem;
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}

	.box.column {
		flex-direction: column;
		align-items: flex-start;
	}

	.toggle {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 1rem;
	}

	.slider {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-top: 0.5rem;
		accent-color: #517188;
	}

	.max-label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		margin-top: 0.5rem;
	}

	input[type='number'] {
		padding: 0.3rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		width: 100px;
	}

	input[type='range'] {
		flex-grow: 1;
	}

	input[type='checkbox'] {
		accent-color: #517188;
		transform: scale(1.5);
	}

	/* Currently unused
	.custom-scraper:hover {
		background-color: #517188;
	}
	*/

	.inputDiv {
		display: flex;
		align-items: center;
		gap: 4rem; /* adds some spacing between text and input */
	}

	.inputDiv input {
		margin-left: auto; /* optional: pushes input to far right */
	}

	#indexerButton {
		padding: 0.65rem 1rem;
		font-size: 1rem;
		font-weight: 500;
	}

	#deleteButton {
		padding: 0.2rem 0.6rem;
		font-size: 1rem;
	}
</style>
