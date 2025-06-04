import { describe, test, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import Page from './+page.svelte';
import Settings from './settings/+page.svelte';

describe('/+page.svelte', () => {
	beforeEach(() => {
		globalThis.fetch = vi.fn();
	});

	test('renders searchbar and search button', () => {
		render(Page);

		expect(screen.getByPlaceholderText('Write your search here!')).toBeInTheDocument();

		expect(screen.getByRole('button', { name: /search/i })).toBeInTheDocument();
	});

	test('check that typing is allowed in searchbar', async () => {
		render(Page);

		const input = screen.getByPlaceholderText('Write your search here!') as HTMLInputElement;

		await fireEvent.input(input, { target: { value: 'a test input' } });

		expect(input.value).toBe('a test input');
	});

	test('shows search results', async () => {
		const mockResults = {
			Query: 'test',
			Results: [
				{ Path: 'local/path/to/file.txt', Score: 0.9, Source: 2 },
				{ Path: 'http://website.com/webpage', Score: 0.79, Source: 1 }
			]
		};

		// mock fetch response
		globalThis.fetch = vi.fn().mockResolvedValueOnce({
			json: async () => mockResults
		});

		render(Page);

		await fireEvent.input(screen.getByPlaceholderText('Write your search here!'), {
			target: { value: 'test' }
		});

		await fireEvent.click(screen.getByRole('button', { name: /search/i }));

		await waitFor(() => {
			expect(screen.getByText('file.txt')).toBeInTheDocument();
			expect(screen.getByRole('button', { name: /download/i })).toBeInTheDocument();
		});
	});

	test('shows a no results found message', async () => {
		const emptyResults = {
			Query: 'test',
			Results: []
		};

		globalThis.fetch = vi.fn().mockResolvedValueOnce({
			json: async () => emptyResults
		});

		render(Page);

		await fireEvent.input(screen.getByPlaceholderText('Write your search here!'), {
			target: { value: 'test' }
		});

		await fireEvent.click(screen.getByRole('button', { name: /search/i }));

		await waitFor(() => {
			expect(screen.getByText('no results for: test')).toBeInTheDocument();
		});
	});

	test('no fetch if search input is empty', async () => {
		render(Page);

		await fireEvent.click(screen.getByRole('button', { name: /search/i }));

		await waitFor(() => {
			expect(fetch).not.toHaveBeenCalled();
		});
	});

	test('search input focused on mount', async () => {
		render(Page);

		const input = screen.getByPlaceholderText('Write your search here!') as HTMLInputElement;
		expect(document.activeElement).toBe(input);
	});

	test('tests file download', async () => {
		const mockResults = {
			Query: 'test',
			Results: [{ Path: 'local/file.pdf', Score: 0.9, Source: 2 }]
		};

		const blob = new Blob(['dummy content']);
		const createObjectURL = vi.fn(() => 'blob:dummy-url');
		const revokeObjectURL = vi.fn();

		globalThis.fetch = vi
			.fn()
			.mockResolvedValueOnce({ json: async () => mockResults }) //search
			.mockResolvedValueOnce({ blob: async () => blob }); //download

		globalThis.URL.createObjectURL = createObjectURL;
		globalThis.URL.revokeObjectURL = revokeObjectURL;

		render(Page);

		await fireEvent.input(screen.getByPlaceholderText('Write your search here!'), {
			target: { value: 'test' }
		});

		await fireEvent.click(screen.getByRole('button', { name: /search/i }));

		const resultHeading = await screen.findByText('file.pdf');
		expect(resultHeading).toBeInTheDocument();

		await fireEvent.click(screen.getByRole('button', { name: /download/i }));

		await waitFor(() => {
			expect(fetch).toHaveBeenCalledWith(
				'http://localhost:8080/download?q=local/file.pdf',
				expect.objectContaining({ method: 'GET' })
			);
			expect(createObjectURL).toHaveBeenCalled();
			expect(revokeObjectURL).toHaveBeenCalled();
		});
	});
});

describe('Settings page', () => {
	test('renders checkboxes  for filters', () => {
		render(Settings);

		expect(screen.getByRole('checkbox', { name: /Files/i })).toBeInTheDocument();
		expect(screen.getByRole('checkbox', { name: /Webpages/i })).toBeInTheDocument();
		expect(screen.getByRole('checkbox', { name: /Show all results/i })).toBeInTheDocument();
	});
});
