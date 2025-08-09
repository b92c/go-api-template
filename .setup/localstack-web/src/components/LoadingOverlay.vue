<template>
  <v-overlay
    v-model="show"
    contained
    persistent
    class="align-center justify-center"
  >
    <div class="text-center">
      <v-progress-circular
        indeterminate
        color="primary"
        size="64"
      ></v-progress-circular>
      <div class="mt-4">
        <h3 class="text-h6 text-white">{{ message }}</h3>
      </div>
    </div>
  </v-overlay>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  show: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: 'Carregando...'
  }
})
</script>
