<template>
  <v-dialog
    eager
    class="max-w-[540px]"
    :model-value="props.show"
    @update:model-value="(val) => emits('update:show', val)"
  >
    <v-card title="Plan your publish">
      <template #text>
        <v-text-field
          clearable
          class="mt-2"
          label="Publish date"
          hint="Your post will hidden for public before this time. Leave blank will publish immediately"
          variant="solo-filled"
          type="datetime-local"
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
const props = defineProps<{ show: boolean; value: string | null }>()
const emits = defineEmits(["update:show", "update:value"])
</script>
