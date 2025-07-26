<template>
  <div class="home-container">
    <h1>Welcome to <i>PianPianino</i></h1>
    <p>Lazyly organize your life, find time for what really <b>counts</b>.</p>
    <p>
      No need to <b>stress</b> click below, to <i>sign up</i> or <i>login</i>.
    </p>
    <div class="button-container">
      <RouterLink to="/register" class="register-link">
        Create Account
      </RouterLink>
      <RouterLink to="/login" class="login-link"> Login </RouterLink>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>
<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.home-container::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: radial-gradient(
      circle at 20% 80%,
      rgba(242, 233, 228, 0.2) 0%,
      transparent 60%
    ),
    radial-gradient(
      circle at 80% 20%,
      rgba(201, 173, 167, 0.15) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 40% 40%,
      rgba(154, 140, 152, 0.08) 0%,
      transparent 70%
    );
  pointer-events: none;
}

.home-container > * {
  position: relative;
  z-index: 1;
}

.button-container {
  display: flex;
  gap: 1.5rem; /* 24px */
  margin-top: 2rem;
  flex-wrap: wrap;
  justify-content: center;
}

.register-link,
.login-link {
  display: inline-block;
  padding: 0.75rem 1.5rem; /* 12px 24px */
  text-decoration: none;
  border-radius: 0.5rem; /* 8px */
  transition: all 0.3s ease;
  font-weight: 600;
  border: 0.125rem solid; /* 2px */
  background-color: white;
}

.register-link {
  border-color: var(--ultra-violet);
  color: var(--ultra-violet);
}

.register-link:hover {
  background-color: var(--ultra-violet);
  color: white;
  transform: translateY(-0.125rem); /* -2px */
  box-shadow: 0 0.25rem 0.75rem rgba(74, 78, 105, 0.3); /* 4px 12px */
}

.login-link {
  border-color: var(--rose-quartz);
  color: var(--rose-quartz);
}

.login-link:hover {
  background-color: var(--rose-quartz);
  color: white;
  transform: translateY(-0.125rem); /* -2px */
  box-shadow: 0 0.25rem 0.75rem rgba(154, 140, 152, 0.3); /* 4px 12px */
}

p {
  font-size: 1.5rem;
  color: var(--space-cadet);
  text-shadow: 0 0.0625rem 0.125rem rgba(255, 255, 255, 0.3); /* 1px 2px */
  text-align: center;
  max-width: 37.5rem; /* 600px */
}

i {
  color: var(--ultra-violet);
  font-style: italic;
}

b {
  color: var(--rose-quartz);
  font-weight: bold;
}

h1 {
  font-size: 2.5rem;
  color: var(--space-cadet);
  text-shadow: 0 0.125rem 0.25rem rgba(255, 255, 255, 0.2); /* 2px 4px */
  margin-bottom: 1.5rem;
  text-align: center;
}

@media (max-width: 48rem) {
  /* 768px */
  .button-container {
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .register-link,
  .login-link {
    width: 12.5rem; /* 200px */
    text-align: center;
  }

  h1 {
    font-size: 2rem;
  }

  p {
    font-size: 1.2rem;
  }
}
</style>
