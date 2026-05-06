<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let reservations = [];
  let error = '';
  let loading = true;

  onMount(async () => {
    if (!$auth.token) { goto('/login'); return; }
    try {
      reservations = await api.get('/reservations');
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  async function cancel(id) {
    if (!confirm('Cancelar esta reserva?')) return;
    try {
      await api.del(`/reservations/${id}`);
      reservations = reservations.filter(r => r.id !== id);
    } catch (e) {
      alert(e.message);
    }
  }

  function statusLabel(s) {
    return { pending: 'Pendente', confirmed: 'Confirmada', refused: 'Recusada' }[s] || s;
  }
</script>

<div class="page">
  <div class="section-header">
    <h2>Minhas Reservas</h2>
    <a href="/reserve" class="btn btn-primary">+ Nova Reserva</a>
  </div>

  {#if loading}
    <p>Carregando...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else if reservations.length === 0}
    <div class="card" style="text-align:center;padding:2rem">
      <p>Você ainda não tem reservas.</p>
      <a href="/reserve" class="btn btn-primary" style="margin-top:1rem;display:inline-block">Fazer uma reserva</a>
    </div>
  {:else}
    <div class="card" style="padding:0;overflow:hidden">
      <table>
        <thead>
          <tr>
            <th>Data</th>
            <th>Horário</th>
            <th>Mesa</th>
            <th>Pessoas</th>
            <th>Status</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each reservations as r}
            <tr>
              <td>{r.date}</td>
              <td>{r.slot_time}</td>
              <td>Mesa {r.table_number} <span style="color:var(--text-muted);font-size:.8rem">({r.location})</span></td>
              <td>{r.party_size}</td>
              <td><span class="badge badge-{r.status}">{statusLabel(r.status)}</span></td>
              <td>
                {#if r.status !== 'refused'}
                  <button class="btn btn-danger" style="padding:.3rem .7rem;font-size:.8rem" on:click={() => cancel(r.id)}>Cancelar</button>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
