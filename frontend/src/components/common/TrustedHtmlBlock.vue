<template>
  <!-- Admin-configured trusted HTML slot. Intentionally no sanitization. -->
  <div v-if="normalizedContent" v-html="normalizedContent"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  content?: string
}>()

const normalizedContent = computed(() => props.content?.trim() || '')
</script>
