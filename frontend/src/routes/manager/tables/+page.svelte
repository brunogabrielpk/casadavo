<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let tables = [];
  let error = '';
  let form = { number: '', location: 'frente', capacity: '', is_active: true };
  let editing = null;

  onMount(async () => {
    if (!$auth.token || $auth.user?.role !== 'gerente') { goto('/login'); return; }
    await load();
  });

  async function load() {
    tables = await api.get('/tables');
  }

  async function save() {
    error = '';
    try {
      const payload = { ...form, number: parseInt(form.number), capacity: parseInt(form.capacity) };
      if (editing) {
        await api.put(`/tables/${editing}`, payload);
        editing = null;
      } else {
        await api.post('/tables', payload);
      }
      form = { number: '', location: 'frente', capacity: '', is_active: true };
      await load();
    } catch (e) {
      error = e.message;
    }
  }

  function startEdit(t) {
    editing = t.id;
    form = { number: t.number, location: t.location, capacity: t.capacity, is_active: t.is_active };
  }

  function cancelEdit() {
    editing = null;
    form = { number: '', location: 'frente', capacity: '', is_active: true };
  }

  async function remove(id) {
    if (!confirm('Remover esta mesa?')) return;
    try {
      await api.del(`/tables/${id}`);
      await load();
    } catch (e) {
      alert(e.message);
    }
  }

  async function toggleActive(t) {
    await api.put(`/tables/${t.id}`, { ...t, is_active: !t.is_active });
    await load();
  }
</script>

<div class="page">
  <h2 style="margin-bottom:1.5rem">Gerenciar Mesas</h2>

  <div class="card" style="margin-bottom:1.5rem">
    <h3 style="margin-bottom:1rem">{editing ? 'Editar Mesa' : 'Nova Mesa'}</h3>
    <form on:submit|preventDefault={save} style="display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:1rem;align-items:end">
      <div class="form-group" style="margin:0">
        <label>Número</label>
        <input type="number" min="1" bind:value={form.number} required />
      </div>
      <div class="form-group" style="margin:0">
        <label>Localização</label>
        <select bind:value={form.location}>
          <option value="frente">Frente</option>
          <option value="fundos">Fundos</option>
        </select>
      </div>
      <div class="form-group" style="margin:0">
        <label>Capacidade</label>
        <input type="number" min="1" bind:value={form.capacity} required />
      </div>
      <div style="display:flex;gap:.5rem">
        <button class="btn btn-primary" type="submit">{editing ? 'Salvar' : 'Adicionar'}</button>
        {#if editing}<button class="btn btn-ghost" type="button" on:click={cancelEdit}>Cancelar</button>{/if}
      </div>
    </form>
    {#if error}<p class="error">{error}</p>{/if}
  </div>

  <div class="table-map">
    {#each tables as t}
      <div class="table-tile {t.is_active ? 'available' : 'inactive'}">
        <div class="tnum">{t.number}</div>
        <div class="tloc">{t.location}</div>
        <div class="tcap">{t.capacity} pessoas</div>
        <div style="margin-top:.5rem;display:flex;gap:.3rem;justify-content:center;flex-wrap:wrap">
          <button class="btn btn-ghost" style="padding:.2rem .5rem;font-size:.75rem" on:click={() => startEdit(t)}>✏️</button>
          <button class="btn" style="padding:.2rem .5rem;font-size:.75rem;background:{t.is_active?'var(--danger)':'var(--success)'};color:#fff" on:click={() => toggleActive(t)}>
            {t.is_active ? 'Bloquear' : 'Ativar'}
          </button>
          <button class="btn btn-danger" style="padding:.2rem .5rem;font-size:.75rem" on:click={() => remove(t.id)}>🗑️</button>
        </div>
      </div>
    {/each}
  </div>
</div>
