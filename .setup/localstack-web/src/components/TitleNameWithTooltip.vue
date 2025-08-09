<template>
  <v-tooltip :text="name" :location="location" :disabled="!showTooltip">
    <template v-slot:activator="{ props }">
      <span 
        class="name-truncate"
        v-bind="props"
      >
        {{ name }}
      </span>
    </template>
  </v-tooltip>
</template>

<script setup>
const props = defineProps({
  name: {
    type: String,
    required: true
  },
  lowercase: {
    type: Boolean,
    default: true
  },
  location: {
    type: String,
    default: 'top'
  },
  customClass: {
    type: String,
    default: ''
  },
  showTooltip: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.name-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: calc(100% - 40px);
}
</style>