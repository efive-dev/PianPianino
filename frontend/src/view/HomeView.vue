<template>
  <HomeIcon></HomeIcon>
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
import HomeIcon from "../components/HomeIcon.vue";
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
  gap: 1.5rem;
  margin-top: 2rem;
  flex-wrap: wrap;
  justify-content: center;
}

.register-link,
.login-link {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  font-weight: 600;
  border: 0.125rem solid;
  background-color: white;
}

.register-link {
  border-color: var(--ultra-violet);
  color: var(--ultra-violet);
}

.register-link:hover {
  background-color: var(--ultra-violet);
  color: white;
  transform: translateY(-0.125rem);
  box-shadow: 0 0.25rem 0.75rem rgba(74, 78, 105, 0.3);
}

.login-link {
  border-color: var(--rose-quartz);
  color: var(--rose-quartz);
}

.login-link:hover {
  background-color: var(--rose-quartz);
  color: white;
  transform: translateY(-0.125rem);
  box-shadow: 0 0.25rem 0.75rem rgba(154, 140, 152, 0.3);
}

@media (max-width: 768px) {
  .button-container {
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .register-link,
  .login-link {
    width: 12.5rem;
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
