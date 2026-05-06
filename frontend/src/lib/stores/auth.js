import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createAuth() {
	const initial = browser
		? { token: localStorage.getItem('token'), user: JSON.parse(localStorage.getItem('user') || 'null') }
		: { token: null, user: null };

	const { subscribe, set, update } = writable(initial);

	return {
		subscribe,
		login(token, user) {
			if (browser) {
				localStorage.setItem('token', token);
				localStorage.setItem('user', JSON.stringify(user));
			}
			set({ token, user });
		},
		logout() {
			if (browser) {
				localStorage.removeItem('token');
				localStorage.removeItem('user');
			}
			set({ token: null, user: null });
		}
	};
}

export const auth = createAuth();
