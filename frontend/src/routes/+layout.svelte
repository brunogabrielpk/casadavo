<script>
  import '../app.css';
  import { auth } from '$lib/stores/auth.js';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  function logout() {
    auth.logout();
    goto('/login');
  }

  $: isManager = $auth.user?.role === 'gerente';
  $: publicRoutes = ['/login', '/register'];
  $: isPublic = publicRoutes.includes($page.url.pathname);
</script>

{#if !isPublic && $auth.token}
<nav>
  <span class="brand">🏠 Casa da Vó</span>
  {#if isManager}
    <a href="/manager">Dashboard</a>
    <a href="/manager/tables">Mesas</a>
    <a href="/manager/availability">Disponibilidade</a>
    <a href="/manager/layout">Layout</a>
  {:else}
    <a href="/reservations">Minhas Reservas</a>
    <a href="/reserve">Nova Reserva</a>
  {/if}
  <span class="spacer"></span>
  <a href="/profile" style="color:#fff;font-size:.9rem">{$auth.user?.name}</a>
  <button class="btn btn-ghost" style="color:#fff;border-color:rgba(255,255,255,.4)" on:click={logout}>Sair</button>
</nav>
{/if}

<slot />
