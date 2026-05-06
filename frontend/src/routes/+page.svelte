<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';

  onMount(() => {
    if (!$auth.token) {
      goto('/login');
    } else if ($auth.user?.role === 'gerente') {
      goto('/manager');
    } else {
      goto('/reservations');
    }
  });
</script>
