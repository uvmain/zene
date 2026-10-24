import { apiKey } from '~/stores/main'

export type QRLoginState = 'pending' | 'scanned' | 'approved' | 'denied' | 'expired' | 'consumed'

export interface QRLoginChallenge {
  state: QRLoginState
  expiresAt: string
}

interface NewQRLoginChallenge {
  phoneSecret: string
  computerSecret: string
  expiresAt: string
}

interface CompletedQRLogin {
  apiKey: string
}

async function qrLoginFetch<T>(path: string, values: Record<string, string> = {}, authenticated = false): Promise<T> {
  const form = new FormData()
  Object.entries(values).forEach(([key, value]) => form.append(key, value))
  if (authenticated) {
    form.append('apiKey', apiKey.value)
    form.append('v', '1.16.1')
    form.append('c', 'zene-frontend')
  }

  const response = await fetch(`/rest/qrlogin/${path}`, {
    method: 'POST',
    body: form,
    cache: 'no-store',
  })
  const data = await response.json() as T & { error?: string }
  if (!response.ok) {
    throw new Error(data.error ?? 'QR login request failed')
  }
  return data
}

export async function createQRLoginChallenge() {
  return qrLoginFetch<NewQRLoginChallenge>('create')
}

export async function getQRLoginChallengeStatus(computerSecret: string) {
  return qrLoginFetch<QRLoginChallenge>('status', { computerSecret })
}

export async function completeQRLoginChallenge(computerSecret: string) {
  return qrLoginFetch<CompletedQRLogin>('complete', { computerSecret })
}

export async function claimQRLoginChallenge(phoneSecret: string) {
  return qrLoginFetch<QRLoginChallenge>('claim', { phoneSecret }, true)
}

export async function decideQRLoginChallenge(phoneSecret: string, decision: 'approve' | 'deny') {
  return qrLoginFetch<QRLoginChallenge>('decide', { phoneSecret, decision }, true)
}
