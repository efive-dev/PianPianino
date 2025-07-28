<template>
  <div class="insert-task-card">
    <n-form label-placement="top" :size="'large'" :label-width="80">
      <n-form-item label="Description" required>
        <n-input
          v-model:value="description"
          placeholder="Task description"
          size="large"
        />
      </n-form-item>

      <n-form-item label="Priority">
        <n-select
          v-model:value="priority"
          :options="priorityOptions"
          placeholder="Select priority"
          size="large"
        />
      </n-form-item>

      <n-form-item>
        <button class="submit-button" @click="handleInsert" :disabled="loading">
          {{ loading ? "Adding..." : "Add Task" }}
        </button>
      </n-form-item>

      <transition name="fade">
        <n-alert v-if="error" type="error" class="form-alert" :bordered="false">
          {{ error }}
        </n-alert>
      </transition>

      <transition name="fade">
        <n-alert
          v-if="success"
          type="success"
          class="form-alert"
          :bordered="false"
        >
          {{ success }}
        </n-alert>
      </transition>
    </n-form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import axios from "axios";

const description = ref("");
const priority = ref("normal");
const loading = ref(false);
const error = ref("");
const success = ref("");

const emit = defineEmits(["task-added"]);

const priorityOptions = [
  { label: "Low", value: "low" },
  { label: "Normal", value: "normal" },
  { label: "High", value: "high" },
];

const handleInsert = async () => {
  error.value = "";
  success.value = "";

  if (!description.value.trim()) {
    error.value = "Description is required";
    return;
  }

  loading.value = true;

  try {
    const token = localStorage.getItem("authToken");
    await axios.post(
      "http://localhost:1323/api/tasks",
      {
        description: description.value,
        priority: priority.value,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      }
    );

    success.value = "Task added!";
    description.value = "";
    priority.value = "normal";
    emit("task-added");
  } catch (err) {
    error.value = err.response?.data?.error || "Failed to add task";
    console.error(err);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.insert-task-card {
  max-width: 32rem;
  margin: 0 auto;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border-radius: 1rem;
  padding: 1.5rem;
  background-color: white;
}

/* Updated button styling to match the rest of the app */
.submit-button {
  width: 100%;
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  font-weight: 600;
  font-size: 1rem;
  border: 0.125rem solid var(--ultra-violet);
  background-color: white;
  color: var(--ultra-violet);
  cursor: pointer;
}

.submit-button:hover:enabled {
  background-color: var(--ultra-violet);
  color: white;
  transform: translateY(-0.125rem);
  box-shadow: 0 0.25rem 0.75rem rgba(74, 78, 105, 0.3);
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-alert {
  margin-top: 1rem;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
