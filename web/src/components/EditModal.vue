<template>
  <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

  <div class="fixed inset-0 z-10 overflow-y-auto">
    <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
      <div
          class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
        <form @submit.prevent="handleSubmit">
          <div>
            <h3 class="text-base font-semibold leading-6 text-gray-900">
              {{ proxy ? 'Edit Proxy' : 'Create New Proxy' }}
            </h3>
            <div class="mt-4 space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Listen URL</label>
                <input
                    v-model="form.listen_url"
                    type="text"
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    required
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Mode</label>
                <select
                    v-model="form.mode"
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                >
                  <option value="path">Path</option>
                  <option value="redirect">Redirect</option>
                  <option value="reverse">Reverse Proxy</option>
                </select>
              </div>

              <!-- Route Condition -->
              <div>
                <label class="block text-sm font-medium text-gray-700">Route Condition</label>
                <div class="mt-2 space-y-4">
                  <div>
                    <select
                        v-model="form.condition.type"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    >
                      <option value="">No condition</option>
                      <option value="header">Header</option>
                      <option value="query">Query Parameter</option>
                      <option value="cookie">Cookie</option>
                      <option value="user_agent">User Agent</option>
                      <option value="language">Language</option>
                    </select>
                  </div>

                  <!-- Parameter Name -->
                  <div v-if="form.condition.type && form.condition.type !== 'language'">
                    <label class="block text-sm font-medium text-gray-700">
                      {{ form.condition.type === 'user_agent' ? 'User Agent Parameter' : 'Parameter Name' }}
                    </label>
                    <select
                        v-if="form.condition.type === 'user_agent'"
                        v-model="form.condition.param_name"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    >
                      <option value="platform">Platform (Mobile/Desktop)</option>
                      <option value="browser">Browser</option>
                    </select>
                    <input
                        v-else
                        v-model="form.condition.param_name"
                        type="text"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        :placeholder="getParamPlaceholder(form.condition.type)"
                    />
                  </div>

                  <!-- Values -->
                  <div v-if="form.condition.type">
                    <label class="block text-sm font-medium text-gray-700">Values</label>
                    <div class="mt-2 space-y-2">
                      <div
                          v-for="(target, index) in form.targets"
                          :key="target.id"
                          class="flex items-center gap-2"
                      >
                        <input
                            v-model="form.condition.values[index]"
                            type="text"
                            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                            :placeholder="getValuePlaceholder(form.condition.type, form.condition.param_name)"
                        />
                        <span class="text-sm text-gray-500">→</span>
                        <span class="text-sm text-gray-700">{{ target.url }}</span>
                      </div>
                    </div>
                    <p class="mt-1 text-sm text-gray-500">
                      {{ getConditionHelp(form.condition.type, form.condition.param_name) }}
                    </p>
                  </div>

                  <!-- Default Target -->
                  <div v-if="form.condition.type">
                    <label class="block text-sm font-medium text-gray-700">Default Target</label>
                    <select
                        v-model="form.condition.default"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    >
                      <option
                          v-for="target in form.targets"
                          :key="target.id"
                          :value="target.id"
                      >
                        {{ target.url }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Tags and Targets -->
              <div>
                <TagsInput
                    v-model="form.tags"
                    :available-tags="availableTags"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Targets</label>
                <div v-for="(target, index) in form.targets" :key="index" class="mt-2 flex items-center gap-2">
                  <input
                      v-model="target.url"
                      type="text"
                      required
                      placeholder="Target URL with optional query"
                      pattern="^(https?:\/\/)?([\da-z.\-]+\.)*[\da-z.\-]+\.([a-z.]{2,6})(\/[\w.\-]*)*(\?[\w.%\-]+(=[\w.%\-]*)?(&[\w.%\-]+(=[\w.%\-]*)?)*)?(#[\w%\-]*)?$"
                      class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  />
                  <input
                      v-model.number="target.weight"
                      type="number"
                      min="0"
                      max="100"
                      placeholder="Weight"
                      class="block w-20 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  />
                  <button
                      type="button"
                      @click="removeTarget(index)"
                      class="text-red-600 hover:text-red-900"
                  >
                    Remove
                  </button>
                </div>
                <button
                    type="button"
                    @click="addTarget"
                    class="mt-2 text-sm text-indigo-600 hover:text-indigo-900"
                >
                  Add Target
                </button>
              </div>
            </div>
          </div>
          <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
            <button
                type="submit"
                class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2"
            >
              {{ proxy ? 'Save Changes' : 'Create' }}
            </button>
            <button
                type="button"
                @click="closeModal"
                class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import TagsInput from "@/components/TagsInput.vue";
import {ref} from "vue";

type Target = {
  id?: string;
  url: string;
  weight: number;
  is_active?: boolean
}

type Form = {
  listen_url: string;
  mode: string;
  tags: string[];
  targets: Target[];
  condition: {
    type: string;
    param_name: string;
    values: string[];
    default: string;
  };
};

type Proxy = Form & {
  id: string
}

const props = defineProps<{
  proxy: Proxy | null;
  availableTags: string[];
  form: Form;
}>()

const emit = defineEmits(['close', 'submit'])

const closeModal = () => {
  emit('close')
}

const form = ref(props.form)

const handleSubmit = () => {
  const formData = {
    listen_url: form.value.listen_url,
    mode: form.value.mode,
    tags: form.value.tags,
    condition: form.value.condition,
    targets: form.value.targets.map(t => ({
      url: t.url,
      weight: t.weight / 100,
      is_active: t.is_active ?? true
    }))
  }
  emit('submit', formData)
}

function addTarget() {
  form.value.targets.push({
    url: '',
    weight: 50,
    is_active: true
  })
}

function removeTarget(index) {
  form.value.targets.splice(index, 1)
}

function getParamPlaceholder(type) {
  switch (type) {
    case 'header':
      return 'X-Custom-Header'
    case 'query':
      return 'version'
    case 'cookie':
      return 'user_preference'
    default:
      return ''
  }
}

function getValuePlaceholder(type, param) {
  switch (type) {
    case 'user_agent':
      return param === 'platform' ? 'mobile/desktop' : 'chrome/firefox/safari/edge/ie'
    case 'language':
      return 'en/es/fr/de'
    default:
      return 'Value for target'
  }
}

function getConditionHelp(type, param) {
  switch (type) {
    case 'user_agent':
      if (param === 'platform') {
        return 'Route traffic based on device type: mobile or desktop'
      }
      return 'Route traffic based on browser type: chrome, firefox, safari, edge, ie, or other'
    case 'language':
      return 'Route traffic based on preferred language from Accept-Language header (e.g., en, es, fr)'
    case 'header':
      return 'Route traffic based on a custom HTTP header value'
    case 'query':
      return 'Route traffic based on a URL query parameter value'
    case 'cookie':
      return 'Route traffic based on a cookie value'
    default:
      return ''
  }
}
</script>