<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';

  let availability = [];
  let allTables = [];
  let tables = [];
  let slots = [];

  let selectedDate = '';
  let selectedAvailID = null;
  let selectedSlotID = '';
  let selectedTableID = '';
  let partySize = 2;
  let notes = '';

  let error = '';
  let success = '';
  let loading = false;
  let loadingSlots = false;

  onMount(async () => {
    if (!$auth.token) { goto('/login'); return; }
    const [av, tb] = await Promise.all([api.get('/availability'), api.get('/tables')]);
    availability = av.filter(a => a.is_open);
    allTables = tb.filter(t => t.is_active);
    tables = allTables;
  });

  async function onDateChange() {
    selectedSlotID = '';
    selectedTableID = '';
    slots = [];
    tables = allTables;
    const av = availability.find(a => a.date === selectedDate);
    if (!av) { selectedAvailID = null; return; }
    selectedAvailID = av.id;
    loadingSlots = true;
    try {
      const [slotData, exclusions] = await Promise.all([
        api.get(`/availability/${av.id}/slots`),
        api.get(`/layout?date=${selectedDate}`)
      ]);
      slots = slotData;
      const excludedIDs = new Set(exclusions.map(e => e.table_id));
      tables = allTables.filter(t => !excludedIDs.has(t.id));
    } finally {
      loadingSlots = false;
    }
  }

  async function submit() {
    error = ''; success = '';
    if (!selectedDate) { error = 'Selecione uma data.'; return; }
    if (!selectedSlotID) { error = 'Selecione um horário.'; return; }
    if (!selectedTableID) { error = 'Selecione uma mesa.'; return; }
    loading = true;
    try {
      await api.post('/reservations', {
        table_id: parseInt(selectedTableID),
        slot_id: parseInt(selectedSlotID),
        date: selectedDate,
        party_size: parseInt(partySize),
        notes
      });
      success = 'Reserva criada! Aguarde a confirmação.';
      selectedDate = ''; selectedSlotID = ''; selectedTableID = ''; notes = ''; partySize = 2; slots = [];
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  $: size = parseInt(partySize) || 1;

  function selectTable(t) {
    if (t.capacity < size) return;
    selectedTableID = t.id;
  }

  $: {
    const selected = tables.find(t => t.id === selectedTableID);
    if (selected && selected.capacity < size) selectedTableID = '';
  }
</script>

<div class="page" style="max-width:700px">
  <h2 style="margin-bottom:1.5rem">Nova Reserva</h2>

  {#if success}
    <div style="background:#d1e7dd;border:1px solid #a3cfbb;padding:1rem;border-radius:var(--radius);margin-bottom:1rem">
      {success} <a href="/reservations">Ver minhas reservas</a>
    </div>
  {/if}

  <div class="card">
    <form on:submit|preventDefault={submit}>

      <!-- Party size first so table filtering is immediate -->
      <div class="form-group">
        <label for="party">Número de pessoas</label>
        <input id="party" type="number" min="1" max="20" bind:value={partySize} required />
      </div>

      <div class="form-group">
        <label for="date">Data</label>
        <select id="date" bind:value={selectedDate} on:change={onDateChange} required>
          <option value="">Selecione uma data</option>
          {#each availability as a}
            <option value={a.date}>{a.date}</option>
          {/each}
        </select>
      </div>

      {#if selectedDate}
        {#if loadingSlots}
          <p style="color:var(--text-muted);margin-bottom:1rem">Carregando horários...</p>
        {:else if slots.length === 0}
          <p style="color:var(--danger);margin-bottom:1rem">Nenhum horário disponível para esta data.</p>
        {:else}
          <div class="form-group">
            <label for="slot">Horário</label>
            <select id="slot" bind:value={selectedSlotID} required>
              <option value="">Selecione um horário</option>
              {#each slots as s}
                <option value={s.id}>{s.slot_time}</option>
              {/each}
            </select>
          </div>
        {/if}

        <div class="form-group">
          <label>Mesa</label>
          <div class="table-map">
            {#each tables as t}
              {@const tooSmall = t.capacity < size}
              <div
                class="table-tile {tooSmall ? 'insufficient' : selectedTableID == t.id ? 'selected' : 'available'}"
                on:click={() => selectTable(t)}
                on:keydown={e => e.key === 'Enter' && selectTable(t)}
                role="button"
                tabindex={tooSmall ? -1 : 0}
                title={tooSmall ? `Capacidade insuficiente (${t.capacity} pessoas)` : ''}
              >
                <div class="tnum">{t.number}</div>
                <div class="tloc">{t.location}</div>
                <div class="tcap">{t.capacity} pessoas</div>
                {#if tooSmall}
                  <div style="font-size:.65rem;color:var(--danger);margin-top:.2rem">Lotada</div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <div class="form-group">
        <label for="notes">Observações (opcional)</label>
        <textarea id="notes" bind:value={notes} rows="2"></textarea>
      </div>

      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn btn-primary" type="submit" disabled={loading || (selectedDate && slots.length === 0)}>
        {loading ? 'Enviando...' : 'Confirmar Reserva'}
      </button>
    </form>
  </div>
</div>

<style>
  .table-tile.insufficient {
    opacity: 0.4;
    filter: grayscale(60%);
    cursor: not-allowed;
    border-color: var(--border);
    background: #f0ebe5;
  }
</style>
