const STORE_NAME = 'episodes'

let db: IDBDatabase

export function createEpisodeStoreIfNotExists() {
  const open = indexedDB.open('data')
  open.onupgradeneeded = () => {
    db = open.result
    if (!db.objectStoreNames.contains(STORE_NAME)) {
      db.createObjectStore(STORE_NAME)
    }
  }
}

export async function getStoredEpisode(key: string): Promise<Blob> {
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

export async function episodeIsStored(key: string): Promise<boolean> {
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

export async function getListOfStoredEpisodes(): Promise<string[]> {
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
      console.log('Fetched list of stored episodes')
    }
  })
}

export function setStoredEpisode(key: string, episode: Blob) {
  const open = indexedDB.open('data')
  open.onsuccess = () => {
    db = open.result
    const transaction = db.transaction(STORE_NAME, 'readwrite')
    const objectStore = transaction.objectStore(STORE_NAME)
    const request = objectStore.put(episode, key)
    request.onerror = () => console.error(request.error)
    transaction.oncomplete = () => db.close()
    console.log(`Stored episode ${key} locally`)
  }
}

export function deleteStoredEpisode(key: string) {
  const open = indexedDB.open('data')
  open.onsuccess = () => {
    db = open.result
    const transaction = db.transaction(STORE_NAME, 'readwrite')
    const objectStore = transaction.objectStore(STORE_NAME)
    const request = objectStore.delete(key)
    request.onerror = () => console.error(request.error)
    transaction.oncomplete = () => db.close()
    console.log(`Deleted episode ${key} from local storage`)
  }
}
