<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let reservations = [];
  let loading = true;
  let error = '';

  onMount(async () => {
    if (!$auth.token || $auth.user?.role !== 'gerente') { goto('/login'); return; }
    try {
      reservations = await api.get('/reservations');
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  async function setStatus(id, status) {
    try {
      await api.put(`/reservations/${id}/status`, { status });
      reservations = reservations.map(r => r.id === id ? { ...r, status } : r);
    } catch (e) {
      alert(e.message);
    }
  }

  function statusLabel(s) {
    return { pending: 'Pendente', confirmed: 'Confirmada', refused: 'Recusada' }[s] || s;
  }

  $: pending = reservations.filter(r => r.status === 'pending');
  $: confirmed = reservations.filter(r => r.status === 'confirmed');
</script>

<div class="page" style="max-width:1100px">
  <h2 style="margin-bottom:1.5rem">Dashboard — Reservas</h2>

  <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;margin-bottom:2rem">
    <div class="card" style="text-align:center">
      <div style="font-size:2rem;font-weight:bold;color:var(--pending)">{pending.length}</div>
      <div style="color:var(--text-muted)">Pendentes</div>
    </div>
    <div class="card" style="text-align:center">
      <div style="font-size:2rem;font-weight:bold;color:var(--success)">{confirmed.length}</div>
      <div style="color:var(--text-muted)">Confirmadas</div>
    </div>
    <div class="card" style="text-align:center">
      <div style="font-size:2rem;font-weight:bold;color:var(--primary)">{reservations.length}</div>
      <div style="color:var(--text-muted)">Total</div>
    </div>
  </div>

  {#if loading}
    <p>Carregando...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else if reservations.length === 0}
    <div class="card" style="text-align:center;padding:2rem">Nenhuma reserva ainda.</div>
  {:else}
    <div class="card" style="padding:0;overflow:hidden">
      <table>
        <thead>
          <tr>
            <th>Data</th>
            <th>Horário</th>
            <th>Mesa</th>
            <th>Cliente</th>
            <th>Telefone</th>
            <th>Pessoas</th>
            <th>Observações</th>
            <th>Status</th>
            <th>Ações</th>
          </tr>
        </thead>
        <tbody>
          {#each reservations as r}
            <tr>
              <td>{r.date}</td>
              <td>{r.slot_time}</td>
              <td>Mesa {r.table_number} <span style="color:var(--text-muted);font-size:.8rem">({r.location})</span></td>
              <td>
                <div>{r.user_name}</div>
                <div style="font-size:.8rem;color:var(--text-muted)">{r.user_email}</div>
              </td>
              <td style="font-size:.85rem">{r.user_phone || '—'}</td>
              <td>{r.party_size}</td>
              <td style="font-size:.85rem;max-width:160px">{r.notes || '—'}</td>
              <td><span class="badge badge-{r.status}">{statusLabel(r.status)}</span></td>
              <td>
                {#if r.status === 'pending'}
                  <div style="display:flex;gap:.4rem">
                    <button class="btn" style="background:var(--success);color:#fff;padding:.3rem .6rem;font-size:.8rem" on:click={() => setStatus(r.id, 'confirmed')}>✓</button>
                    <button class="btn btn-danger" style="padding:.3rem .6rem;font-size:.8rem" on:click={() => setStatus(r.id, 'refused')}>✗</button>
                  </div>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
