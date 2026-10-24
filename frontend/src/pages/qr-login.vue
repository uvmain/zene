<script setup lang="ts">
import type { QRLoginChallenge } from '~/logic/qrLogin'
import { claimQRLoginChallenge, decideQRLoginChallenge } from '~/logic/qrLogin'

const router = useRouter()
const phoneSecret = ref('')
const challenge = ref<QRLoginChallenge | null>(null)
const error = ref('')
const submitting = ref(false)

function secretFromFragment() {
  return new URLSearchParams(window.location.hash.slice(1)).get('challenge') || ''
}

onMounted(async () => {
  phoneSecret.value = secretFromFragment()
  if (!phoneSecret.value) {
    error.value = 'This QR login link is invalid.'
    return
  }
  try {
    challenge.value = await claimQRLoginChallenge(phoneSecret.value)
    if (!['pending', 'scanned'].includes(challenge.value.state)) {
      error.value = 'This QR login request is no longer available.'
    }
  }
  catch (cause) {
    error.value = cause instanceof Error
      ? cause.message
      : 'This QR login request is unavailable.'
  }
})

async function decide(decision: 'approve' | 'deny') {
  if (!phoneSecret.value || submitting.value) {
    return
  }
  submitting.value = true
  error.value = ''
  try {
    challenge.value = await decideQRLoginChallenge(phoneSecret.value, decision)
  }
  catch (cause) {
    error.value = cause instanceof Error
      ? cause.message
      : 'Could not update this QR login request.'
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="my-auto p-4 background-1 flex flex-col gap-6 max-w-screen items-center justify-center">
    <div class="text-muted text-center flex flex-col gap-5 max-w-400px items-center">
      <div class="text-xl font-bold">
        Approve QR login
      </div>
      <template v-if="challenge?.state === 'approved'">
        <p>This browser login has been approved. You can return to the music player.</p>
        <ZButton @click="router.push('/')">
          Continue
        </ZButton>
      </template>
      <template v-else-if="challenge?.state === 'denied'">
        <p>This browser login request was denied.</p>
      </template>
      <template v-else-if="!error">
        <p>Approving will log your account into another browser. Only continue if you started this login request.</p>
        <div class="flex gap-3">
          <ZButton :disabled="submitting" @click="decide('approve')">
            {{ submitting ? 'Working…' : 'Approve login' }}
          </ZButton>
          <ZButton :disabled="submitting" @click="decide('deny')">
            Deny
          </ZButton>
        </div>
      </template>
      <p v-if="error" class="text-red-500">
        {{ error }}
      </p>
    </div>
  </div>
</template>
