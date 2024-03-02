<template>
  <v-dialog v-model="editor.show" class="max-w-[540px]">
    <v-card title="New post">
      <v-form>
        <v-card-text>
          <v-textarea
            required
            hide-details
            variant="outlined"
            label="What's happened?!"
          />

          <div class="flex mt-1">
            <v-tooltip text="Planned publish" location="start">
              <template #activator="{ props }">
                <v-btn v-bind="props" type="button" variant="text" icon="mdi-calendar" size="small" />
              </template>
            </v-tooltip>
            <v-tooltip text="Categories" location="start">
              <template #activator="{ props }">
                <v-btn v-bind="props" type="button" variant="text" icon="mdi-shape" size="small" />
              </template>
            </v-tooltip>
            <v-tooltip text="Media" location="start">
              <template #activator="{ props }">
                <v-btn v-bind="props" type="button" variant="text" icon="mdi-camera" size="small" />
              </template>
            </v-tooltip>
          </div>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn type="reset" color="grey" @click="editor.show = false">Cancel</v-btn>
          <v-btn type="submit" @click.prevent>Publish</v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { useEditor } from "@/stores/editor";

const editor = useEditor();
</script>
