<template>
  <div>
    <label :for="id" class="block text-sm font-medium text-gray-700">{{ label }}</label>
    <div class="mt-1">
      <div class="flex flex-wrap gap-2 p-2 border rounded-md">
        <!-- Existing tags -->
        <span
          v-for="tag in modelValue"
          :key="tag"
          class="inline-flex items-center px-2 py-1 rounded-md text-sm font-medium bg-indigo-100 text-indigo-700"
        >
          {{ tag }}
          <button
            type="button"
            @click="removeTag(tag)"
            class="ml-1 inline-flex text-indigo-400 hover:text-indigo-500"
          >
            <span class="sr-only">Remove tag</span>
            ×
          </button>
        </span>

        <!-- Input for new tag -->
        <input
          :id="id"
          v-model="newTag"
          type="text"
          @keydown.enter.prevent="addTag"
          @keydown.tab.prevent="addTag"
          @keydown.comma.prevent="addTag"
          placeholder="Add tags..."
          class="flex-1 min-w-[120px] border-0 bg-transparent p-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm"
        />
      </div>
    </div>

    <!-- Available tags suggestions -->
    <div v-if="suggestions.length > 0" class="mt-1">
      <div class="text-sm text-gray-500">Available tags:</div>
      <div class="flex flex-wrap gap-1 mt-1">
        <button
          v-for="tag in suggestions"
          :key="tag"
          @click="addExistingTag(tag)"
          type="button"
          class="inline-flex items-center px-2 py-1 rounded-md text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200"
        >
          {{ tag }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true
  },
  label: {
    type: String,
    default: 'Tags'
  },
  id: {
    type: String,
    default: 'tags-input'
  },
  availableTags: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const newTag = ref('')

const suggestions = computed(() => {
  return props.availableTags.filter(tag => 
    !props.modelValue.includes(tag) &&
    (!newTag.value || tag.toLowerCase().includes(newTag.value.toLowerCase()))
  )
})

function addTag() {
  const tag = newTag.value.trim()
  if (tag && !props.modelValue.includes(tag)) {
    emit('update:modelValue', [...props.modelValue, tag])
  }
  newTag.value = ''
}

function addExistingTag(tag) {
  if (!props.modelValue.includes(tag)) {
    emit('update:modelValue', [...props.modelValue, tag])
  }
}

function removeTag(tagToRemove) {
  emit('update:modelValue', props.modelValue.filter(tag => tag !== tagToRemove))
}
</script>
