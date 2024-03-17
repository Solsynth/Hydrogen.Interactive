<template>
  <v-list density="comfortable">
    <v-list-subheader>Realms</v-list-subheader>
    <v-list-item
      v-for="item in realms"
      exact
      prepend-icon="mdi-account-multiple"
      :to="{ name: 'realms.details', params: { id: item.id } }"
      :title="item.name"
    />

    <v-divider v-if="realms.length > 0" class="border-opacity-75 my-2" />

    <v-list-item
      prepend-icon="mdi-plus"
      title="Create a realm"
      :disabled="!id.userinfo.isLoggedIn"
      @click="creating = true"
    />
  </v-list>

  <v-dialog v-model="creating" class="max-w-[540px]">
    <v-card title="Create a realm" prepend-icon="mdi-account-multiple-plus" :loading="loading">
      <v-form @submit.prevent="submit">
        <v-card-text>
          <v-text-field
            label="Name"
            variant="outlined"
            density="comfortable"
            v-model="requestData.name"
          />
          <v-textarea
            label="Description"
            variant="outlined"
            density="comfortable"
            v-model="requestData.description"
          />
          <v-select
            label="Realm type"
            item-title="label"
            item-value="value"
            variant="outlined"
            density="comfortable"
            :items="realmTypeOptions"
            v-model="requestData.realm_type"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn type="reset" color="grey-darken-3" @click="creating = false">Cancel</v-btn>
          <v-btn type="submit" :disabled="loading">Save</v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>

  <!-- @vue-ignore -->
  <v-snackbar v-model="error" :timeout="5000">Something went wrong... {{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { getAtk, useUserinfo } from "@/stores/userinfo";
import { useEditor } from "@/stores/editor";

const id = useUserinfo();
const editor = useEditor();

const realms = computed(() => editor.availableRealms);
const requestData = ref({
  name: "",
  description: "",
  realm_type: 0
});

const realmTypeOptions = [
  { label: "Public Realm", value: 0 },
  { label: "Restricted Realm", value: 1 },
  { label: "Private Realm", value: 2 }
];

const creating = ref(false);

const error = ref<string | null>(null);
const reverting = ref(false);
const loading = ref(false);

async function list() {
  reverting.value = true;
  try {
    await editor.listRealms();
  } catch (err) {
    error.value = (err as Error).message;
  }
  reverting.value = false;
}

async function submit(evt: SubmitEvent) {
  const form = evt.target as HTMLFormElement;
  const payload = requestData.value;
  if (!payload.name) return;

  loading.value = true;
  const res = await fetch("/api/realms", {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getAtk()}` },
    body: JSON.stringify(payload)
  });
  if (res.status !== 200) {
    error.value = await res.text();
  } else {
    await list();
    form.reset();
    creating.value = false;
  }
  loading.value = false;
}
</script>
