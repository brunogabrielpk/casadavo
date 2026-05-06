<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let tables = [];
  let exclusions = [];
  let date = new Date().toISOString().slice(0, 10);
  let loading = false;
  let error = '';

  $: excludedIDs = new Map(exclusions.map(e => [e.table_id, e.id]));

  onMount(async () => {
    if (!$auth.token || $auth.user?.role !== 'gerente') { goto('/login'); return; }
    tables = (await api.get('/tables')).filter(t => t.is_active);
    await loadExclusions();
  });

  async function loadExclusions() {
    loading = true;
    error = '';
    try {
      exclusions = await api.get(`/layout?date=${date}`);
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function toggle(table) {
    const exclusionID = excludedIDs.get(table.id);
    try {
      if (exclusionID) {
        await api.del(`/layout/${exclusionID}`);
        exclusions = exclusions.filter(e => e.id !== exclusionID);
      } else {
        const ex = await api.post('/layout', { table_id: table.id, date });
        exclusions = [...exclusions, ex];
      }
    } catch (e) {
      error = e.message;
    }
  }
</script>

<div class="page">
  <div class="section-header">
    <h2>Layout de Mesas</h2>
  </div>

  <div class="card" style="margin-bottom:1.5rem;display:flex;gap:1rem;align-items:center;flex-wrap:wrap">
    <div class="form-group" style="margin:0;flex:1;min-width:200px">
      <label for="date">Data</label>
      <input id="date" type="date" bind:value={date} on:change={loadExclusions} />
    </div>
    <p style="color:var(--text-muted);font-size:.85rem;margin:0">
      Mesas marcadas como <strong>indisponíveis</strong> não aparecerão para reservas nesta data.
    </p>
  </div>

  {#if error}<p class="error">{error}</p>{/if}

  {#if loading}
    <p>Carregando...</p>
  {:else}
    <div class="card" style="padding:0;overflow:hidden">
      <table>
        <thead>
          <tr>
            <th>Mesa</th>
            <th>Local</th>
            <th>Capacidade</th>
            <th>Status nesta data</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each tables as t}
            {@const excluded = excludedIDs.has(t.id)}
            <tr>
              <td>Mesa {t.number}</td>
              <td>{t.location}</td>
              <td>{t.capacity} pessoas</td>
              <td>
                <span class="badge" class:badge-confirmed={!excluded} class:badge-refused={excluded}>
                  {excluded ? 'Indisponível' : 'Disponível'}
                </span>
              </td>
              <td>
                <button
                  class="btn"
                  class:btn-danger={!excluded}
                  class:btn-primary={excluded}
                  style="padding:.3rem .7rem;font-size:.8rem"
                  on:click={() => toggle(t)}
                >
                  {excluded ? 'Disponibilizar' : 'Bloquear'}
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
