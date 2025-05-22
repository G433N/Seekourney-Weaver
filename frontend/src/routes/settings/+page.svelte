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
	} from '../../lib/stores/settings';

	import { get } from 'svelte/store';

	import { onMount } from 'svelte';

	interface IndexerResult {
		ID: number;
		Name: string;
		ExecPath: string;
		Args: string[];
		Port: number;
	}

	let indexerInput: '';
	// let submittedIndexer: string = '';
	//let indexerList: IndexerResult[] = $state([]);

	async function fetchIndexers(): Promise<void> {
		try {
			const res = await fetch('http://localhost:8080/all/indexers');
			const data: IndexerResult[] = await res.json();
			indexerList.set(data);
		} catch (err) {
			console.error('Failed to fetch indexers: ', err);
		}
	}

	async function addIndexer(): Promise<void> {
		if (!indexerInput) return;

		try {
			await fetch('http://localhost:8080/push/indexer', {
				method: 'POST',
				headers: {
					'Content-Type': 'text/plain'
				},
				body: indexerInput
			});

			indexerInput = '';
			await fetchIndexers();
		} catch (err) {
			console.error('Failed to fetch indexer: ', err);
		}
	}

	/*
	async function deleteIndexer(indexer: IndexerResult): Promise<void> {
		//fetch(`http://localhost:8080/addIndexer?q=${indexer}`); //TODO: what name??

		//indexerList = indexerList.filter((elem) => elem.Id !== indexer.Id);
		console.log('delete indexer');
		// TODO: unsure if we get a response?
	}
	*/

	onMount(() => {
		fetchIndexers();
	});

</script>

<main style="max-width: 600px;">
	<h2>Filter</h2>
	<div class="box">
		<h3>Show:</h3>
		<label class="toggle">
			<input type="checkbox" bind:checked={$showFiles} />
			Files
		</label>
		<label class="toggle">
			<input type="checkbox" bind:checked={$showWebpages} />
			Webpages
		</label>
	</div>

	<div class="box column">
		<label class="toggle">
			<input type="checkbox" bind:checked={$showAllResults} />
			Show all results
		</label>
		
		<label class="max-label" class:disabled-label={$showAllResults}>
			<div class="inputDiv">
				Max results shown:
				<input type="number" bind:value={$maxResults} disabled={$showAllResults} />
			</div>
		</label>
	</div>

	<h2>
		CPU usage
		<span class="tooltip-wrapper">
			<span class="info-icon">?</span>
			<span class="tooltip-text">
				CPU usage determines how much of the CPU the system is allowed to use. More cores means 
				faster searching but may make other processes running on your device slower.
			</span>
		</span>
	</h2>
	<div class="box column">
		<label class="toggle">
			<input type="checkbox" bind:checked={$cpuDefault} />
			default CPU usage
		</label>
		<div class="slider">
			<label for="cpuSlider" class:disabled-label={$cpuDefault}>CPU cores used:</label>
			<input id="cpuSlider" type="range" min="1" max={$maxCores} bind:value={$cpuCores} disabled={$cpuDefault}/>
			<span class:disabled-label={$cpuDefault}>{$cpuCores}</span>
		</div>
	</div>

	<h2>
		Indexer
		<span class="tooltip-wrapper">
			<span class="info-icon">?</span>
			<span class="tooltip-text">
				An indexer is a small program that looks through files or websites and creates a list of what they contain
				, like words in documents, so you can search and find things quickly.
				<br>
				<br>
				How to use:
				<br>
				1. Enter the absolute path in the search field, the indexer must be located among your files.
				<br>
				2. Click "Add" - this will register you indexer and start using it automatically
				<br>
				<br>
				You can add multiple indexers to handle different kinds of data.
			</span>
		</span>
	</h2>
	<div class="box column">
		<label class="max-label">
			<p>
				You can add your own custom indexer by providing the absolute path to a local executable on your device.
			</p>
			<div>
				Indexer path:
				<input id="InputIndexer" type="text" bind:value={indexerInput} placeholder="/user/example/custom-indexer" />
				<button id="indexerButton" onclick={() => addIndexer()}> Add </button>
			</div>
		</label>

		{#if $indexerList.length > 0}
			{#each $indexerList as indexer}
				<div class="inputDiv">
					<h3>
						{indexer.Name}, ID: {indexer.ID}
					</h3>
					<small>
						({indexer.ExecPath}) | {indexer.Args.join(' ')}
					</small>
					<!--
					<button id="deleteButton" onclick={() => deleteIndexer(indexer)}> Delete </button>
					-->
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

	.disabled-label {
		color: #999;
		opacity: 0.7;
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

	p {
		font-weight: 400;
	}

	.tooltip-wrapper {
		position: relative;
  		display: inline-flex;
  		align-items: center;
  		justify-content: center;
  		cursor: help;
  		width: 1.2em;
  		height: 1.2em;
	}

	.info-icon {
	  	font-size: 0.7em;
  		width: 1.1em;
  		height: 1.1em;
  		border-radius: 50%;
  		background-color: white;
  		border: 2px solid black;
  		color: black;
  		font-weight: bold;
  		display: flex;
  		align-items: center;
  		justify-content: center;
  		line-height: 1;
	}

	.tooltip-text {
		visibility: hidden;
		width: 320px;
		background-color: #333;
		color: #fff;
		text-align: left;
		border-radius: 8px;
		padding: 0.75rem;
		position: absolute;
		z-index: 1;
		bottom: 125%;
		left: 50%;
		transform: translateX(-50%);
		opacity: 0;
		transition: opacity 0.2s;
		font-size: 0.9rem;
		line-height: 1.4;
		pointer-events: none;
	}

	.tooltip-wrapper:hover .tooltip-text {
		visibility: visible;
		opacity: 1;
	}

	#InputIndexer {
		width: 350px;
	}

</style>
