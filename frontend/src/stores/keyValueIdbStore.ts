import { debugLog } from '~/logic/logger'

const STORE_NAME = 'kv'

let db: IDBDatabase

export function createKVStoreIfNotExists() {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('kv')
    open.onupgradeneeded = () => {
      db = open.result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME)
      }
    }
  }
}

export async function getStoredKV(key: string): Promise<string | null> {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('kv')
    return new Promise<string | null>((resolve, reject) => {
      open.onsuccess = () => {
        db = open.result
        const transaction = db.transaction(STORE_NAME)
        const objectStore = transaction.objectStore(STORE_NAME)
        const request = objectStore.get(key)
        request.onerror = () => reject(request.error)
        request.onsuccess = () => resolve(request.result as string | null)
        transaction.oncomplete = () => db.close()
      }
    })
  }
  else {
    return Promise.reject(new Error('indexedDB is not available'))
  }
}

export function setStoredKV(key: string, value: string) {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('kv')
    open.onsuccess = () => {
      db = open.result
      const transaction = db.transaction(STORE_NAME, 'readwrite')
      const objectStore = transaction.objectStore(STORE_NAME)
      const request = objectStore.put(value, key)
      request.onerror = () => console.error(request.error)
      transaction.oncomplete = () => db.close()
      debugLog(`Stored value ${key} locally`)
    }
  }
}

export function deleteStoredKV(key: string) {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('kv')
    open.onsuccess = () => {
      db = open.result
      const transaction = db.transaction(STORE_NAME, 'readwrite')
      const objectStore = transaction.objectStore(STORE_NAME)
      const request = objectStore.delete(key)
      request.onerror = () => console.error(request.error)
      transaction.oncomplete = () => db.close()
      debugLog(`Deleted value ${key} from local storage`)
    }
  }
}
