<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-base font-semibold leading-6 text-gray-900">Proxies</h1>
        <p class="mt-2 text-sm text-gray-700">A list of all proxies in your A/B testing system.</p>
      </div>
      <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
        <button
            @click="openCreateModal"
            class="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        >
          Add proxy
        </button>
      </div>
    </div>

    <!-- Filter by tags -->
    <div class="mt-4">
      <TagsInput
          v-model="selectedTags"
          label="Filter by tags"
          :available-tags="availableTags"
          @update:modelValue="filterProxies"
      />
    </div>

    <!-- Proxies List -->
    <div class="mt-8 flow-root">
      <ProxiesList :filteredProxies="filteredProxies" :sortBy="sortBy" :sortDesc="sortDesc" @sort="handleSort"
                   @delete="deleteProxy" @edit="editProxy" @viewHistory="viewHistory"/>
    </div>

    <!-- Pagination -->
    <Pagination :currentPage="currentPage" :itemsPerPage="itemsPerPage" @changePage="changePage"
                :totalProxies="totalProxies"/>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="relative z-10">
      <EditModal :proxy="editingProxy" :form="form" :available-tags="availableTags" @close="closeModal" @submit="handleSubmit"/>
    </div>

    <!-- History Modal -->
    <div v-if="showHistoryModal" class="relative z-10">
      <HistoryModal :proxy-id="selectedProxyId" @close="closeHistoryModal"/>
    </div>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import axios from 'axios'
import TagsInput from '@/components/TagsInput.vue'
import Pagination from "@/components/Pagination.vue";
import ProxiesList from "@/components/ProxiesList.vue";
import EditModal from "@/components/EditModal.vue";
import HistoryModal from "@/components/HistoryModal.vue";

const proxies = ref([])
const totalProxies = ref(0)
const currentPage = ref(1)
const itemsPerPage = ref(10)
const sortBy = ref('id')
const sortDesc = ref(false)
const selectedTags = ref([])
const availableTags = ref([])
const showModal = ref(false)
const showHistoryModal = ref(false)
const editingProxy = ref(null)
const selectedProxyId = ref(null)

const form = ref({
  listen_url: '',
  mode: 'reverse',
  tags: [],
  targets: [],
  condition: {
    type: '',
    param_name: '',
    values: [],
    default: ''
  },
})

const filteredProxies = computed(() => {
  if (selectedTags.value.length === 0) return proxies.value
  return proxies.value.filter(proxy =>
      selectedTags.value.every(tag => proxy.tags?.includes(tag))
  )
})

async function loadProxies() {
  try {
    const response = await axios.get('/api/proxies', {
      params: {
        limit: itemsPerPage.value,
        offset: (currentPage.value - 1) * itemsPerPage.value,
        sortBy: sortBy.value,
        sortDesc: sortDesc.value
      }
    })
    proxies.value = response.data.items.map((proxy) => {
      return {
        ...proxy,
        condition: {
          ...proxy.condition,
          values: proxy.condition?.values ? proxy.targets.map(({id}) => proxy.condition.values[id]) : []
        }
      }
    })
    totalProxies.value = response.data.total
  } catch (error) {
    console.error('Failed to load proxies:', error)
  }
}

function handleSort(column) {
  if (sortBy.value === column) {
    sortDesc.value = !sortDesc.value
  } else {
    sortBy.value = column
    sortDesc.value = false
  }
  loadProxies()
}

function changePage(page) {
  currentPage.value = page
  loadProxies()
}

async function filterProxies() {
  if (selectedTags.value.length === 0) {
    await loadProxies()
  } else {
    try {
      const response = await axios.get('/api/proxies/by-tags', {
        params: {tags: selectedTags.value.join(',')}
      })
      proxies.value = response.data.proxies || []
    } catch (error) {
      console.error('Failed to filter proxies:', error)
    }
  }
}

function openCreateModal() {
  editingProxy.value = null
  form.value = {
    listen_url: '',
    mode: 'reverse',
    tags: [],
    targets: [],
    condition: {
      type: '',
      param_name: '',
      values: [],
      default: ''
    }
  }
  showModal.value = true
}

function editProxy(proxy) {
  editingProxy.value = proxy
  form.value = {
    listen_url: proxy.listen_url,
    mode: proxy.mode,
    tags: proxy.tags || [],
    targets: proxy.targets.map(t => ({
      url: t.url,
      weight: t.weight * 100,
      is_active: t.is_active
    })),
    condition: proxy.condition || {
      type: '',
      param_name: '',
      values: [],
      default: ''
    }
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingProxy.value = null
  form.value = {
    listen_url: '',
    mode: 'reverse',
    tags: [],
    targets: [],
    condition: {
      type: '',
      param_name: '',
      values: [],
      default: ''
    },
  }
}

function viewHistory(proxy) {
  selectedProxyId.value = proxy.id
  showHistoryModal.value = true
}

function closeHistoryModal() {
  showHistoryModal.value = false
  selectedProxyId.value = null
}

async function handleSubmit(formData) {
  try {
    if (editingProxy.value) {
      await Promise.all([
        axios.put(`/api/proxies/${editingProxy.value.id}/targets`, formData),
        axios.put(`/api/proxies/${editingProxy.value.id}/tags`, {tags: formData.tags})
      ])
    } else {
      await axios.post('/api/proxies', formData)
    }

    await loadProxies()
    closeModal()
  } catch (error) {
    console.error('Failed to save proxy:', error)
    alert(error.response?.data?.error || 'Failed to save proxy')
  }
}

async function deleteProxy(id) {
  if (!confirm('Are you sure you want to delete this proxy?')) return

  try {
    await axios.delete(`/api/proxies/${id}`)
    await loadProxies()
  } catch (error) {
    console.error('Failed to delete proxy:', error)
  }
}

onMounted(async () => {
  await loadProxies()
  const tagsResponse = await axios.get('/api/tags')
  availableTags.value = tagsResponse.data.tags || []
})
</script>
