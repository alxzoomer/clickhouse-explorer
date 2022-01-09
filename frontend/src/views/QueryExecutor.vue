<template>
  <div class="q-pa-md">
    <q-card>
      <q-card-section>
        <v-ace-editor
          v-model="editor"
          v-model:value="query"
          @init="editorInit"
          lang="sql"
          theme="chrome"
          :wrap="true"
          :min-lines="5"
          :max-lines="20"
        />
      </q-card-section>
      <q-separator />
      <q-card-actions>
        <q-btn flat @click="execQuery" class="text-green" hint="test">Execute</q-btn>
      </q-card-actions>
    </q-card>
    <q-separator spaced dark />
    <q-table
      title="Query result"
      :rows="dbrows"
      :columns="dbcols"
      :loading="loading"
      :rows-per-page-options="[25, 50, 100]"
    >
      <template v-slot:body-cell="props">
        <q-td :props="props">{{ formatCell(props) }}</q-td>
      </template>
      <template v-slot:no-data="{ message }">
        <div class="col-1">
          <q-icon name="warning" class="q-table__bottom-nodata-icon text-red" v-if="errorMessage" />
          <q-icon name="info" class="q-table__bottom-nodata-icon" v-else />
        </div>
        <div class="col-11 error">{{ queryStatus(message) }}</div>
      </template>
    </q-table>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { VAceEditor } from 'vue3-ace-editor';
import type { Ace } from 'ace-builds';

import 'ace-builds/src-noconflict/mode-sql';
import 'ace-builds/src-noconflict/theme-chrome';

export default defineComponent({
  components: {
    VAceEditor,
  },
  setup() {
    const dbrows = ref(new Array<unknown>());
    const dbcols = ref(new Array<any>());
    const loading = ref(false);
    const query = ref('');
    const editor = ref({});
    const errorMessage = ref('');

    function formatCell(props: any): unknown {
      const { value } = props;
      return Array.isArray(value) ? JSON.stringify(value) : value;
    }
    async function execQuery() {
      const url = '/api/v1/query';
      // const data = { query: 'select * from test.example_table order by 1' };
      const data = { query: query.value };
      loading.value = true;
      errorMessage.value = '';
      try {
        const resp = await fetch(url, {
          method: 'POST',
          body: JSON.stringify(data),
          headers: {
            'Content-Type': 'application/json',
          },
        });
        if (!resp.ok) {
          throw await resp.json();
        }
        const json = await resp.json();
        const jrows = json.rows;
        if (jrows.length > 0) {
          const cols = Object.keys(jrows[0]).map((k) => ({ name: k, label: k, field: k }));
          dbcols.value.splice(0, dbcols.value.length, ...cols);
        } else {
          dbcols.value.splice(0, dbcols.value.length);
        }
        dbrows.value.splice(0, dbrows.value.length, ...jrows);
      } catch (e: any) {
        errorMessage.value = e.error;
      } finally {
        loading.value = false;
      }
    }
    function editorInit(v: Ace.Editor) {
      query.value = '-- comment\nselect * from test.example_table;';
      v.commands.addCommand({
        name: 'execQuery',
        bindKey: { win: 'Ctrl-Enter', mac: 'Command-Enter' },
        exec() {
          execQuery();
        },
      });
    }
    function queryStatus(defaultMessage: string): string {
      return errorMessage.value || defaultMessage;
    }
    return {
      dbrows,
      dbcols,
      loading,
      formatCell,
      execQuery,
      query,
      editorInit,
      editor,
      queryStatus,
      errorMessage,
    };
  },
});
</script>

<style lang="sass" scoped>
div.error
  overflow-wrap: break-word
</style>
