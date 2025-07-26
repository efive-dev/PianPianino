<template>
  <n-space vertical size="large" class="dashboard-container">
    <h1>Dashboard</h1>

    <InsertTask @task-added="fetchTasks" />

    <n-card
      class="tasks-section"
      title="Your Tasks"
      :bordered="false"
      size="medium"
    >
      <div v-if="loading" class="loading">Loading tasks...</div>
      <div v-else-if="tasks.length === 0" class="no-tasks">No tasks found.</div>

      <n-space vertical size="large">
        <n-card
          v-for="task in tasks"
          :key="task.id"
          size="small"
          class="task-card"
          :hoverable="true"
          :style="{ backgroundColor: cardColor(task.priority) }"
        >
          <div class="task-item">
            <div class="task-text">
              <div :class="{ completed: task.completed }" class="description">
                {{ task.description }}
              </div>
              <n-tag
                :type="priorityType(task.priority)"
                size="small"
                class="priority-tag"
              >
                {{ capitalize(task.priority) }}
              </n-tag>
            </div>

            <n-button
              size="small"
              tertiary
              circle
              @click="handleDelete(task.id)"
            >
              X
            </n-button>
          </div>
        </n-card>
      </n-space>
    </n-card>
  </n-space>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { NCard, NTag, NButton, NSpace } from "naive-ui";
import InsertTask from "../../components/InsertTask.vue";
import axios from "axios";

const tasks = ref([]);
const loading = ref(false);

const fetchTasks = async () => {
  loading.value = true;
  try {
    const token = localStorage.getItem("authToken");
    const response = await axios.get("http://localhost:1323/api/tasks", {
      headers: { Authorization: `Bearer ${token}` },
    });
    tasks.value = response.data.tasks || [];
  } catch {
    tasks.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchTasks);

const capitalize = (str) =>
  str ? str.charAt(0).toUpperCase() + str.slice(1) : "";

const priorityType = (priority) => {
  switch (priority) {
    case "high":
      return "error";
    case "normal":
      return "warning";
    case "low":
      return "success";
    default:
      return "default";
  }
};

// Assign light pastel backgrounds based on priority
const cardColor = (priority) => {
  switch (priority) {
    case "high":
      return "#fff1f0"; // light red
    case "normal":
      return "#fffbe6"; // light yellow
    case "low":
      return "#f6ffed"; // light green
    default:
      return "#f8f9fa"; // neutral light gray
  }
};

const handleDelete = async (id) => {
  loading.value = true;

  try {
    const token = localStorage.getItem("authToken");
    await axios.delete(`http://localhost:1323/api/tasks/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    tasks.value = tasks.value.filter((task) => task.id !== id);
  } catch (err) {
    const error = ref("");
    error.value = err.response?.data?.error || "Failed to delete task";
    console.error(err);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.dashboard-container {
  max-width: 48rem;
  margin: 0 auto;
  padding: 2rem 1rem;
}

h1 {
  font-size: 2.5rem;
  text-align: center;
  margin-bottom: 2rem;
}

.tasks-section {
  background-color: white;
  border-radius: 0.75rem;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.loading,
.no-tasks {
  font-size: 1.25rem;
  color: #5c6ac4;
  text-align: center;
  margin: 1rem 0;
}

.task-card {
  border-radius: 0.75rem;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  transition: all 0.25s ease;
}

.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-text {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  max-width: 75%;
}

.description {
  font-weight: 600;
  font-size: 1.1rem;
  word-break: break-word;
}

.description.completed {
  color: #aaa;
  text-decoration: line-through;
}

.priority-tag {
  width: fit-content;
}
</style>
