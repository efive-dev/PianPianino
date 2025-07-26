<template>
  <div class="login-container">
    <n-card title="Login" class="login-card" hoverable>
      <n-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-placement="top"
        size="large"
      >
        <n-form-item label="Username" path="username">
          <n-input
            v-model:value="form.username"
            placeholder="Enter your username"
          />
        </n-form-item>
        <n-form-item label="Password" path="password">
          <n-input
            type="password"
            v-model:value="form.password"
            placeholder="Enter a secure password"
            show-password-on="click"
          />
        </n-form-item>
        <n-button type="primary" block @click="handleLogin" :loading="loading">
          Login
        </n-button>
        <n-alert
          v-if="errorMessage"
          type="error"
          class="error-alert"
          :bordered="false"
        >
          {{ errorMessage }}
        </n-alert>
        <n-alert
          v-if="successMessage"
          type="success"
          class="success-alert"
          :bordered="false"
        >
          {{ successMessage }}
        </n-alert>
      </n-form>
    </n-card>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { NForm, NFormItem, NInput, NButton, NCard, NAlert } from "naive-ui";
import axios from "axios";

const router = useRouter();

const form = ref({
  username: "",
  password: "",
});

const rules = {
  username: {
    required: true,
    message: "Username is required",
    trigger: ["blur", "input"],
  },
  password: {
    required: true,
    message: "Password is required",
    trigger: ["blur", "input"],
  },
};

const formRef = ref(null);
const errorMessage = ref("");
const successMessage = ref("");
const loading = ref(false);

const handleLogin = async () => {
  errorMessage.value = "";
  successMessage.value = "";

  const valid = await formRef.value?.validate();
  if (!valid) return;

  loading.value = true;

  try {
    const response = await axios.post(
      "http://localhost:1323/login",
      form.value
    );

    // Store the JWT token
    if (response.data.token) {
      localStorage.setItem("authToken", response.data.token);

      // Set default authorization header for future requests
      axios.defaults.headers.common[
        "Authorization"
      ] = `Bearer ${response.data.token}`;
    }

    successMessage.value = response.data.message;
    form.value.username = "";
    form.value.password = "";

    setTimeout(() => {
      router.push("/dashboard");
    }, 1500);
  } catch (err) {
    errorMessage.value = err.response?.data?.error || "Login failed";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-attachment: fixed;
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.login-container::before {
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

.login-card {
  width: 100%;
  max-width: 25rem;
  border-radius: 1rem;
  box-shadow: 0 0.25rem 1.25rem rgba(34, 34, 59, 0.1);
  background-color: white;
  border: 0.125rem solid var(--pale-dogwood);
  position: relative;
  z-index: 1;
}

.login-card :deep(.n-card-header) {
  color: var(--space-cadet);
  font-weight: 600;
  font-size: 1.5rem;
  text-align: center;
}

.login-card :deep(.n-form-item-label) {
  color: var(--ultra-violet);
  font-weight: 500;
}

.login-card :deep(.n-input) {
  border-color: var(--rose-quartz);
  border-radius: 0.5rem;
}

.login-card :deep(.n-input:hover) {
  border-color: var(--ultra-violet);
}

.login-card :deep(.n-input:focus-within) {
  border-color: var(--ultra-violet);
  box-shadow: 0 0 0 0.125rem rgba(74, 78, 105, 0.2);
}

.login-card :deep(.n-button) {
  background-color: var(--ultra-violet);
  color: white;
  border-radius: 0.5rem;
  border: none;
  font-weight: 600;
  transition: all 0.3s ease;
  margin-top: 1rem;
}

.login-card :deep(.n-button:hover) {
  background-color: var(--space-cadet);
  transform: translateY(-0.0625rem);
  box-shadow: 0 0.25rem 0.75rem rgba(74, 78, 105, 0.3);
}

.login-card :deep(.n-button:active) {
  transform: translateY(0);
}

.login-card :deep(.n-button--loading) {
  background-color: var(--rose-quartz);
}

.error-alert,
.success-alert {
  margin-top: 1rem;
  border-radius: 0.5rem;
}

.login-card :deep(.n-alert--error-type) {
  background-color: rgba(193, 18, 31, 0.1);
  border: 0.0625rem solid rgba(193, 18, 31, 0.3);
  color: #c1121f;
}

.login-card :deep(.n-alert--success-type) {
  background-color: rgba(74, 78, 105, 0.1);
  border: 0.0625rem solid rgba(74, 78, 105, 0.3);
  color: var(--ultra-violet);
}

@media (max-width: 768px) {
  .login-container {
    padding: 1rem;
  }

  .login-card {
    max-width: 100%;
  }
}
</style>
