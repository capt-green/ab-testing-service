<template>
  <div class="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6 mt-4">
    <div class="flex flex-1 justify-between sm:hidden">
      <button
          @click="changePage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
      >
        Previous
      </button>
      <button
          @click="changePage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage === totalPages }"
      >
        Next
      </button>
    </div>
    <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-gray-700">
          Showing
          <span class="font-medium">{{ startIndex + 1 }}</span>
          to
          <span class="font-medium">{{ endIndex }}</span>
          of
          <span class="font-medium">{{ totalProxies }}</span>
          results
        </p>
      </div>
      <div>
        <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
          <button
              v-for="page in totalPages"
              :key="page"
              @click="changePage(page)"
              class="relative inline-flex items-center px-4 py-2 text-sm font-semibold"
              :class="[
                page === currentPage
                  ? 'z-10 bg-indigo-600 text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
                  : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:outline-offset-0',
                page === 1 ? 'rounded-l-md' : '',
                page === totalPages ? 'rounded-r-md' : ''
              ]"
          >
            {{ page }}
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed } from 'vue';

const props = defineProps({
  currentPage: {
    type: Number,
    required: true,
  },
  totalProxies: {
    type: Number,
    required: true,
  },
  itemsPerPage: {
    type: Number,
    default: 10,
  },
});

const emit = defineEmits(['changePage']);

const totalPages = computed(() => Math.ceil(props.totalProxies / props.itemsPerPage))
const startIndex = computed(() => (props.currentPage - 1) * props.itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + props.itemsPerPage, props.totalProxies))

const changePage = (page) => {
  emit('changePage', page)
}

</script>