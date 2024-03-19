<template>
  <v-dialog
    eager
    class="max-w-[540px]"
    :model-value="props.show"
    @update:model-value="(val) => emits('update:show', val)"
  >
    <v-card title="Change your audiences">
      <template #text>
        <v-select
          clearable
          class="mt-2"
          label="Realm"
          hint="This field will only show realms you joined. Leave blank to publish this post in public area."
          variant="solo-filled"
          item-title="name"
          item-value="id"
          :items="realms.available"
          :model-value="props.value"
          @update:model-value="(val) => emits('update:value', val)"
        />
      </template>
      <template #actions>
        <v-btn class="ms-auto" text="Ok" @click="emits('update:show', false)"></v-btn>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor";
import { useRealms } from "@/stores/realms";

const realms = useRealms();
const editor = useEditor();

const props = defineProps<{ show: boolean; value: string | null }>();
const emits = defineEmits(["update:show", "update:value"]);
</script>
