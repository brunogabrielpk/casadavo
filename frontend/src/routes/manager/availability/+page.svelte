<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let availability = [];
  let selectedAv = null;
  let slots = [];
  let newDate = '';
  let newSlotTime = '';
  let error = '';

  onMount(async () => {
    if (!$auth.token || $auth.user?.role !== 'gerente') { goto('/login'); return; }
    await loadAv();
  });

  async function loadAv() {
    availability = await api.get('/availability');
  }

  async function addDay() {
    error = '';
    try {
      await api.post('/availability', { date: newDate, is_open: true, auto_confirm: false });
      newDate = '';
      await loadAv();
    } catch (e) {
      error = e.message;
    }
  }

  async function toggleOpen(a) {
    await api.put(`/availability/${a.id}`, { ...a, is_open: !a.is_open });
    await loadAv();
  }

  async function toggleAutoConfirm(a) {
    await api.put(`/availability/${a.id}`, { ...a, auto_confirm: !a.auto_confirm });
    await loadAv();
  }

  async function selectDay(a) {
    selectedAv = a;
    slots = await api.get(`/availability/${a.id}/slots`);
  }

  async function addSlot() {
    if (!newSlotTime) return;
    try {
      await api.post(`/availability/${selectedAv.id}/slots`, { slot_time: newSlotTime });
      newSlotTime = '';
      slots = await api.get(`/availability/${selectedAv.id}/slots`);
    } catch (e) {
      alert(e.message);
    }
  }

  async function removeSlot(id) {
    await api.del(`/slots/${id}`);
    slots = slots.filter(s => s.id !== id);
  }
</script>

<div class="page" style="max-width:900px">
  <h2 style="margin-bottom:1.5rem">Disponibilidade</h2>

  <div style="display:grid;grid-template-columns:1fr 1fr;gap:1.5rem;align-items:start">
    <!-- Left: day list -->
    <div>
      <div class="card" style="margin-bottom:1rem">
        <h4 style="margin-bottom:.75rem">Abrir novo dia</h4>
        <div style="display:flex;gap:.5rem">
          <input type="date" bind:value={newDate} />
          <button class="btn btn-primary" on:click={addDay}>Adicionar</button>
        </div>
        {#if error}<p class="error">{error}</p>{/if}
      </div>

      <div class="card" style="padding:0;overflow:hidden">
        {#if availability.length === 0}
          <p style="padding:1rem;color:var(--text-muted)">Nenhum dia configurado.</p>
        {:else}
          {#each availability as a}
            <div
              style="padding:.75rem 1rem;border-bottom:1px solid var(--border);cursor:pointer;background:{selectedAv?.id===a.id?'#fff8ee':'transparent'}"
              on:click={() => selectDay(a)}
              on:keydown={e => e.key==='Enter' && selectDay(a)}
              role="button"
              tabindex="0"
            >
              <div style="display:flex;align-items:center;justify-content:space-between">
                <strong>{a.date}</strong>
                <span class="badge {a.is_open ? 'badge-confirmed' : 'badge-refused'}">{a.is_open ? 'Aberto' : 'Fechado'}</span>
              </div>
              <div style="display:flex;gap:.5rem;margin-top:.5rem;flex-wrap:wrap">
                <button class="btn btn-ghost" style="padding:.2rem .6rem;font-size:.75rem" on:click|stopPropagation={() => toggleOpen(a)}>
                  {a.is_open ? 'Fechar' : 'Abrir'}
                </button>
                <button class="btn btn-ghost" style="padding:.2rem .6rem;font-size:.75rem" on:click|stopPropagation={() => toggleAutoConfirm(a)}>
                  Confirm. auto: {a.auto_confirm ? '✓' : '✗'}
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Right: slots for selected day -->
    <div>
      {#if selectedAv}
        <div class="card">
          <h4 style="margin-bottom:.75rem">Horários — {selectedAv.date}</h4>
          <div style="display:flex;gap:.5rem;margin-bottom:1rem">
            <input type="time" bind:value={newSlotTime} />
            <button class="btn btn-primary" on:click={addSlot}>Adicionar</button>
          </div>
          {#if slots.length === 0}
            <p style="color:var(--text-muted)">Nenhum horário definido.</p>
          {:else}
            <ul style="list-style:none;display:flex;flex-direction:column;gap:.5rem">
              {#each slots as s}
                <li style="display:flex;align-items:center;justify-content:space-between;padding:.4rem .6rem;background:#f9f4ef;border-radius:var(--radius)">
                  <span style="font-size:1.05rem">{s.slot_time}</span>
                  <button class="btn btn-danger" style="padding:.2rem .5rem;font-size:.8rem" on:click={() => removeSlot(s.id)}>✗</button>
                </li>
              {/each}
            </ul>
          {/if}
        </div>
      {:else}
        <div class="card" style="text-align:center;padding:2rem;color:var(--text-muted)">
          Selecione um dia para gerenciar horários.
        </div>
      {/if}
    </div>
  </div>
</div>
