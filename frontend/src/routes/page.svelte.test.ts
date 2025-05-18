import { describe, test, expect, vi, beforeEach } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import Page from './+page.svelte';

describe('/+page.svelte', () => {
	beforeEach(() => {
		// mock global fetch before every test to not depend on backend behaviour for unit tests
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
				{ Path: 'document', Score: 0.9, Source: 1 },
				{ Path: 'website', Score: 0.79, Source: 2 }
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

		// may need to change depending on how we display results
		await waitFor(() => {
			expect(screen.getByText('document')).toBeInTheDocument();
			expect(screen.getByText('website')).toBeInTheDocument();
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
});
