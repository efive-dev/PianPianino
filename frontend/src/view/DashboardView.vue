<template>
  <HomeIcon></HomeIcon>
  <div class="dashboard-page">
    <div class="dashboard-header">
      <h1>Dashboard</h1>
      <button class="logout-button" @click="handleLogout">Logout</button>
    </div>

    <n-space vertical size="large" class="dashboard-container">
      <div class="button-container">
        <button class="add-task-btn" @click="showInsert = !showInsert">
          {{ showInsert ? "Hide Task Form" : "Add Task" }}
        </button>
      </div>

      <transition name="fade-slide">
        <div v-if="showInsert" class="insert-task-wrapper">
          <InsertTask @task-added="onTaskAdded" />
        </div>
      </transition>

      <n-card
        class="tasks-section"
        title="Your Tasks"
        :bordered="false"
        size="medium"
      >
        <n-space
          justify="space-between"
          align="center"
          style="margin-bottom: 1rem"
        >
          <div style="font-weight: bold">Your Tasks</div>
          <div style="display: block; gap: 0.5rem">
            <n-select
              v-model:value="sortOrder"
              :options="[
                { label: 'Newest First', value: 'desc' },
                { label: 'Oldest First', value: 'asc' },
              ]"
              size="small"
              style="width: 10rem; margin-bottom: 0.5rem"
            />
            <n-select
              v-model:value="priorityFilter"
              :options="[
                { label: 'All', value: 'all' },
                { label: 'High', value: 'high' },
                { label: 'Normal', value: 'normal' },
                { label: 'Low', value: 'low' },
              ]"
              size="small"
              style="width: 10rem"
            />
          </div>
        </n-space>

        <div v-if="tasks.length === 0" class="no-tasks">No tasks found.</div>

        <n-space vertical size="large">
          <n-card
            v-for="task in sortedTasks"
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
                <div class="created-at">{{ formatDate(task.created_at) }}</div>
                <n-tag
                  :type="priorityType(task.priority)"
                  size="small"
                  class="priority-tag"
                >
                  {{ capitalize(task.priority) }}
                </n-tag>
              </div>

              <div class="task-actions">
                <n-button size="small" tertiary @click="handleDelete(task.id)">
                  Delete
                </n-button>
                <n-button size="small" secondary @click="toggleTask(task.id)">
                  {{ task.completed ? "Undo" : "Complete" }}
                </n-button>
              </div>
            </div>
          </n-card>
        </n-space>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import InsertTask from "../components/InsertTask.vue";
import HomeIcon from "../components/HomeIcon.vue";

const router = useRouter();
const tasks = ref([]);
const loading = ref(false);
const sortOrder = ref("desc");
const showInsert = ref(true);
const priorityFilter = ref("all");

const handleLogout = () => {
  localStorage.removeItem("authToken");
  router.push("/");
};

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

const onTaskAdded = () => {
  fetchTasks();
};

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

const cardColor = (priority) => {
  switch (priority) {
    case "high":
      return "#fff1f0";
    case "normal":
      return "#fffbe6";
    case "low":
      return "#f6ffed";
    default:
      return "#f8f9fa";
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
    console.error(err.response?.data?.error || "Failed to delete task");
  } finally {
    loading.value = false;
  }
};

const toggleTask = async (id) => {
  loading.value = true;
  try {
    const token = localStorage.getItem("authToken");
    await axios.patch(
      `http://localhost:1323/api/tasks/${id}/toggle`,
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      }
    );

    const index = tasks.value.findIndex((task) => task.id === id);
    if (index !== -1) {
      tasks.value[index].completed = !tasks.value[index].completed;
    }
  } catch (err) {
    console.error("Failed to toggle task:", err);
  } finally {
    loading.value = false;
  }
};

const formatDate = (isoString) => {
  const date = new Date(isoString);
  return date.toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  });
};

const sortedTasks = computed(() => {
  let filtered = [...tasks.value];

  if (priorityFilter.value !== "all") {
    filtered = filtered.filter((task) => {
      return task.priority === priorityFilter.value;
    });
  }

  return filtered.sort((a, b) => {
    const timeA = new Date(a.created_at).getTime();
    const timeB = new Date(b.created_at).getTime();
    return sortOrder.value === "asc" ? timeA - timeB : timeB - timeA;
  });
});
</script>

<style scoped>
.dashboard-page {
  position: relative;
}

.dashboard-header {
  position: relative;
  padding: 1rem 1rem 0;
  text-align: center;
}

.dashboard-header h1 {
  font-size: 2rem;
  font-weight: 700;
  margin: 0;
}

.logout-button {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.5rem 1rem;
  font-weight: 600;
  background-color: white;
  color: var(--rose-quartz);
  border: 2px solid var(--rose-quartz);
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-button:hover {
  background-color: var(--rose-quartz);
  color: white;
  transform: translateY(-0.125rem);
  box-shadow: 0 0.25rem 0.75rem rgba(154, 140, 152, 0.3);
}

.dashboard-container {
  max-width: 48rem;
  margin: 0 auto;
  padding: 1rem;
}

.button-container {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.add-task-btn {
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  font-weight: 600;
  border: 0.125rem solid var(--ultra-violet);
  background-color: white;
  color: var(--ultra-violet);
  cursor: pointer;
}

.add-task-btn:hover {
  background-color: var(--ultra-violet);
  color: white;
  transform: translateY(-0.125rem);
  box-shadow: 0 0.25rem 0.75rem rgba(74, 78, 105, 0.3);
}

.insert-task-wrapper {
  margin-bottom: 1.5rem;
}

/* Animation */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
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

.created-at {
  font-size: 0.85rem;
  color: #888;
}

.priority-tag {
  width: fit-content;
}

.task-actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}
</style>
