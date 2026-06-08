<script setup lang="ts">
import { ref, watch, nextTick } from "vue"
import { useAuthStore } from "@/stores/auth"
import { useSettingsStore } from "@/stores/settings"

const auth = useAuthStore()
const settings = useSettingsStore()
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ "update:modelValue": [value: boolean] }>()

const isRegister = ref(false)
const loginUsername = ref("")
const loginPassword = ref("")
const loginConfirm = ref("")
const loginError = ref("")
const showPwd = ref(false)



async function handleAuth() {
  loginError.value = ""
  if (isRegister.value) {
    if (loginPassword.value !== loginConfirm.value) { loginError.value = "�������벻һ��"; return }
    const err = await auth.register(loginUsername.value, loginPassword.value)
    if (err) { loginError.value = err; return }
  } else {
    const err = await auth.login(loginUsername.value, loginPassword.value)
    if (err) { loginError.value = err; return }
  }
  await nextTick()
  resetForm()
}

function resetForm() {
  emit("update:modelValue", false)
  loginUsername.value = ""
  loginPassword.value = ""
  loginConfirm.value = ""
  loginError.value = ""
}
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)"
    max-width="400" persistent transition="dialog-bottom-transition">
    <v-card class="rounded-xl" rounded="xl">
      <v-card-title class="d-flex align-center pa-4 pb-0">
        <div class="d-flex align-center ga-2">
          <v-icon color="primary">mdi-account-lock</v-icon>
          <span class="text-h6 font-weight-medium">{{ isRegister ? 'ע���˺�' : '��¼' }}</span>
        </div>
        <v-spacer />
        <v-btn icon="mdi-close" size="small" variant="text" @click.stop="emit('update:modelValue', false)" />
      </v-card-title>
      <v-card-text class="pa-4">
        <v-alert v-if="loginError" :text="loginError" type="error" variant="tonal" density="compact" closable class="mb-4 rounded-lg" @click:close="loginError = ''" />
        <v-form @submit.prevent="handleAuth">
          <v-text-field v-model="loginUsername" label="�û���" variant="outlined" density="comfortable" hide-details class="mb-3" prepend-inner-icon="mdi-account-outline" />
          <v-text-field v-model="loginPassword" :type="showPwd ? 'text' : 'password'" label="����" variant="outlined" density="comfortable" hide-details class="mb-3" prepend-inner-icon="mdi-lock-outline">
            <template #append-inner>
              <v-btn :icon="showPwd ? 'mdi-eye-off' : 'mdi-eye'" size="x-small" variant="text" @click.stop="showPwd = !showPwd" />
            </template>
          </v-text-field>
          <v-text-field v-if="isRegister" v-model="loginConfirm" :type="showPwd ? 'text' : 'password'" label="ȷ������" variant="outlined" density="comfortable" hide-details class="mb-3" prepend-inner-icon="mdi-lock-outline" />
          <v-btn type="submit" color="primary" variant="flat" size="large" block class="rounded-pill mt-2"
            :disabled="isRegister && !settings.allowRegister">{{ isRegister && !allowRegister ? 'ע���ѹر�' : (isRegister ? 'ע�Ტ��¼' : '��¼') }}</v-btn>
        </v-form>
      </v-card-text>
      <v-card-actions class="pa-4 pt-0 d-flex justify-center">
        <v-btn v-if="settings.allowRegister" variant="text" size="small" class="text-caption text-medium-emphasis"
          @click.stop="isRegister = !isRegister; loginError = ''">{{ isRegister ? '�����˺ţ�ȥ��¼' : 'û���˺ţ�ȥע��' }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

