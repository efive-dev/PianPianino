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
        <n-button
          type="primary"
          block
          :loading="loading"
          size="large"
          @click="handleInsert"
        >
          {{ loading ? "Adding..." : "Add Task" }}
        </n-button>
      </n-form-item>

      <n-alert v-if="error" type="error" class="mt-4" :bordered="false">
        {{ error }}
      </n-alert>
      <n-alert v-if="success" type="success" class="mt-4" :bordered="false">
        {{ success }}
      </n-alert>
    </n-form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useMessage } from "naive-ui";
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
  console.log("Insert triggered");

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

.success-msg {
  color: #52c41a;
  font-size: 0.95rem;
  margin-top: 1rem;
  text-align: center;
}
</style>
