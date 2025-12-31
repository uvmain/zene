import { useDebug } from '~/composables/useDebug'

const { debugLog } = useDebug()

const STORE_NAME = 'episodes'

let db: IDBDatabase

export function createEpisodeStoreIfNotExists() {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    open.onupgradeneeded = () => {
      db = open.result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME)
      }
    }
  }
}

export async function getStoredEpisode(key: string): Promise<Blob> {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    return new Promise<Blob>((resolve, reject) => {
      open.onsuccess = () => {
        db = open.result
        const transaction = db.transaction(STORE_NAME)
        const objectStore = transaction.objectStore(STORE_NAME)
        const request = objectStore.get(key)
        request.onerror = () => reject(request.error)
        request.onsuccess = () => resolve(request.result as Blob)
        transaction.oncomplete = () => db.close()
      }
    })
  }
  else {
    return Promise.reject(new Error('indexedDB is not available'))
  }
}

export async function episodeIsStored(key: string): Promise<boolean> {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    return new Promise<boolean>((resolve, reject) => {
      open.onsuccess = () => {
        db = open.result
        const transaction = db.transaction(STORE_NAME)
        const objectStore = transaction.objectStore(STORE_NAME)
        const request = objectStore.get(key)
        request.onerror = () => reject(request.error)
        request.onsuccess = () => resolve(request.result !== undefined)
        transaction.oncomplete = () => db.close()
      }
    })
  }
  else {
    return Promise.reject(new Error('indexedDB is not available'))
  }
}

export async function getListOfStoredEpisodes(): Promise<string[]> {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    return new Promise<string[]>((resolve, reject) => {
      open.onsuccess = () => {
        db = open.result
        const transaction = db.transaction(STORE_NAME)
        const objectStore = transaction.objectStore(STORE_NAME)
        const request = objectStore.getAllKeys()
        request.onerror = () => reject(request.error)
        request.onsuccess = () => resolve(request.result as string[])
        transaction.oncomplete = () => db.close()
        debugLog('Fetched list of stored episodes')
      }
    })
  }
  else {
    return Promise.reject(new Error('indexedDB is not available'))
  }
}

export function setStoredEpisode(key: string, episode: Blob) {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    open.onsuccess = () => {
      db = open.result
      const transaction = db.transaction(STORE_NAME, 'readwrite')
      const objectStore = transaction.objectStore(STORE_NAME)
      const request = objectStore.put(episode, key)
      request.onerror = () => console.error(request.error)
      transaction.oncomplete = () => db.close()
      debugLog(`Stored episode ${key} locally`)
    }
  }
}

export function deleteStoredEpisode(key: string) {
  if (typeof window !== 'undefined' && 'indexedDB' in window) {
    const open = indexedDB.open('data')
    open.onsuccess = () => {
      db = open.result
      const transaction = db.transaction(STORE_NAME, 'readwrite')
      const objectStore = transaction.objectStore(STORE_NAME)
      const request = objectStore.delete(key)
      request.onerror = () => console.error(request.error)
      transaction.oncomplete = () => db.close()
      debugLog(`Deleted episode ${key} from local storage`)
    }
  }
}
