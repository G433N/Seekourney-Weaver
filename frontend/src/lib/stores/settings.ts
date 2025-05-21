import { writable } from '@svelte/store';

export const showFiles = writable(true);
export const showWebpages = writable(true);
export const showAllResults = writable(true);
export const maxResults = writable(100);

export const cpuDefault = writable(true);
export const maxCores = writable(8);
export const cpuCores = writable(4);

export const indexerList = writable<{ Name: string; Id: number; Port: number }[]>([]);