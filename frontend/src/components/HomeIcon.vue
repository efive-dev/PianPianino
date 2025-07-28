<template>
  <RouterLink to="/" class="icon-wrapper">
    <img src="/icon.png" alt="PianPianino Logo" />
  </RouterLink>
</template>

<script setup>
import { RouterLink } from 'vue-router';
</script>

<style scoped>
.icon-wrapper {
  position: fixed;
  top: 1rem;
  left: 1rem;
  background-color: var(--isabelline);
  border-radius: 50%;
  padding: 0.5rem;
  box-shadow: 0 0.25rem 1rem rgba(34, 34, 59, 0.15);
  width: 4rem;
  height: 4rem;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  transition: transform 0.2s ease;
}

.icon-wrapper:hover {
  transform: scale(1.05);
}

.icon-wrapper img {
  max-width: 100%;
  height: auto;
  border-radius: 50%;
}

@media (max-width: 768px) {
  .icon-wrapper {
    width: 3rem;
    height: 3rem;
    top: 0.75rem;
    left: 0.75rem;
  }
}
</style>
