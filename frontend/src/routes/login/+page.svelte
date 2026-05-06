<script>
  import { api } from '$lib/api.js';
  import { auth } from '$lib/stores/auth.js';
  import { goto } from '$app/navigation';

  let email = '';
  let password = '';
  let error = '';
  let loading = false;

  async function submit() {
    error = '';
    loading = true;
    try {
      const data = await api.post('/auth/login', { email, password });
      auth.login(data.token, data.user);
      if (data.user.role === 'gerente') goto('/manager');
      else goto('/reservations');
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }
</script>

<div style="display:flex;align-items:center;justify-content:center;min-height:100vh;background:var(--bg)">
  <div class="card" style="width:100%;max-width:400px">
    <h2 style="text-align:center;margin-bottom:1.5rem;color:var(--primary)">🏠 Casa da Vó</h2>
    <form on:submit|preventDefault={submit}>
      <div class="form-group">
        <label for="email">E-mail</label>
        <input id="email" type="email" bind:value={email} required />
      </div>
      <div class="form-group">
        <label for="password">Senha</label>
        <input id="password" type="password" bind:value={password} required />
      </div>
      {#if error}<p class="error">{error}</p>{/if}
      <button class="btn btn-primary" style="width:100%;margin-top:.5rem" type="submit" disabled={loading}>
        {loading ? 'Entrando...' : 'Entrar'}
      </button>
    </form>
    <p style="text-align:center;margin-top:1rem;font-size:.9rem">
      Não tem conta? <a href="/register">Cadastrar</a>
    </p>
  </div>
</div>
