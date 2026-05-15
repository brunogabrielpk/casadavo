<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let phone = '';
  let success = '';
  let error = '';
  let loading = false;

  onMount(() => {
    if (!$auth.token) { goto('/login'); return; }
    phone = $auth.user?.phone ?? '';
  });

  async function save() {
    success = ''; error = '';
    loading = true;
    try {
      const updated = await api.put('/auth/me', { phone });
      auth.login($auth.token, updated);
      success = 'Telefone atualizado com sucesso.';
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }
</script>

<div style="display:flex;align-items:center;justify-content:center;min-height:80vh;background:var(--bg)">
  <div class="card" style="width:100%;max-width:420px">
    <h2 style="margin-bottom:1.5rem;color:var(--primary)">Meu Perfil</h2>

    <div class="form-group">
      <label>Nome</label>
      <input type="text" value={$auth.user?.name ?? ''} disabled style="opacity:.6;cursor:not-allowed" />
    </div>
    <div class="form-group">
      <label>E-mail</label>
      <input type="email" value={$auth.user?.email ?? ''} disabled style="opacity:.6;cursor:not-allowed" />
    </div>

    <form on:submit|preventDefault={save}>
      <div class="form-group">
        <label for="phone">Telefone</label>
        <input id="phone" type="tel" bind:value={phone} placeholder="(00) 00000-0000" />
      </div>

      {#if success}<p style="color:var(--success);margin-bottom:.75rem">{success}</p>{/if}
      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn btn-primary" style="width:100%" type="submit" disabled={loading}>
        {loading ? 'Salvando...' : 'Salvar'}
      </button>
    </form>
  </div>
</div>
