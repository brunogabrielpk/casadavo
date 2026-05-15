import { auth } from './stores/auth.js';
import { get } from 'svelte/store';

const BASE = '/api';

function headers() {
	const { token } = get(auth);
	return {
		'Content-Type': 'application/json',
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};
}

async function request(method, path, body) {
	const res = await fetch(BASE + path, {
		method,
		headers: headers(),
		body: body !== undefined ? JSON.stringify(body) : undefined
	});
	if (res.status === 204) return null;
	if (res.status === 401) {
		auth.logout();
		window.location.href = '/login';
		return;
	}
	const data = await res.json();
	if (!res.ok) throw new Error(data.error || res.statusText);
	return data;
}

export const api = {
	get: (path) => request('GET', path),
	post: (path, body) => request('POST', path, body),
	put: (path, body) => request('PUT', path, body),
	del: (path) => request('DELETE', path)
};
