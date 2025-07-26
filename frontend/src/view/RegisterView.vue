<template>
  <div class="register-container">
    <n-card title="Create an Account" class="register-card" hoverable>
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
        <n-button
          type="primary"
          block
          @click="handleRegister"
          :loading="loading"
        >
          Register
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
import { NForm, NFormItem, NInput, NButton, NCard, NAlert } from "naive-ui";
import axios from "axios";

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

const handleRegister = async () => {
  errorMessage.value = "";
  successMessage.value = "";

  const valid = await formRef.value?.validate();
  if (!valid) return;

  loading.value = true;

  try {
    const response = await axios.post(
      "http://localhost:1323/register",
      form.value
    );
    successMessage.value = response.data.message;
    form.value.username = "";
    form.value.password = "";
  } catch (err) {
    errorMessage.value = err.response?.data?.error || "Registration failed";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f9fafb;
  padding: 2rem;
}

.register-card {
  width: 100%;
  max-width: 400px;
  border-radius: 16px;
  box-shadow: 0 4px 14px rgba(82, 53, 123, 0.1);
  background-color: #ffffff;
  border: 1px solid #b2d8ce;
}

.register-card :deep(.n-button) {
  background-color: #52357b;
  color: white;
  border-radius: 8px;
  transition: background-color 0.3s ease;
}

.register-card :deep(.n-button:hover) {
  background-color: #5459ac;
}

.error-alert,
.success-alert {
  margin-top: 1rem;
}
</style>
