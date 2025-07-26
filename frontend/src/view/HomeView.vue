<template>
  <div class="home-container">
    <h1>Welcome to the App</h1>
    <RouterLink to="/register" class="register-link">
      Create Account
    </RouterLink>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #f9fafb;
  padding: 2rem;
}

.register-link {
  display: inline-block;
  padding: 12px 24px;
  background-color: #52357b;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  transition: background-color 0.3s ease;
  margin-top: 1rem;
}

.register-link:hover {
  background-color: #5459ac;
}
</style>