import { debugLog } from '~/logic/logger'

class GenericIndexedDbStore {
  private dbName = 'genericData'
  private storeName = 'genericStore'
  private db: IDBDatabase | null = null
  private dbPromise: Promise<IDBDatabase> | null = null

  async openDb(): Promise<IDBDatabase> {
    if (this.db) {
      return Promise.resolve(this.db)
    }

    if (this.dbPromise) {
      return this.dbPromise
    }

    this.dbPromise = new Promise<IDBDatabase>((resolve, reject) => {
      const probe = indexedDB.open(this.dbName)

      probe.onsuccess = () => {
        const probeDb = probe.result
        const needsUpgrade = !probeDb.objectStoreNames.contains(this.storeName)
        const version = probeDb.version
        probeDb.close()

        if (!needsUpgrade) {
          const request = indexedDB.open(this.dbName, version)
          request.onsuccess = () => {
            this.db = request.result
            this.dbPromise = null
            resolve(this.db)
          }
          request.onerror = () => {
            this.dbPromise = null
            reject(request.error)
          }
        }
        else {
          const request = indexedDB.open(this.dbName, version + 1)
          request.onupgradeneeded = () => {
            const db = request.result
            if (!db.objectStoreNames.contains(this.storeName)) {
              db.createObjectStore(this.storeName)
            }
          }
          request.onsuccess = () => {
            this.db = request.result
            this.dbPromise = null
            resolve(this.db)
          }
          request.onerror = () => {
            this.dbPromise = null
            reject(request.error)
          }
          request.onblocked = () => {
            debugLog(`IndexedDB blocked while opening store: ${this.storeName}`)
          }
        }
      }

      probe.onerror = () => {
        this.dbPromise = null
        reject(probe.error)
      }
    })

    debugLog('GenericIndexedDbStore opened')
    return this.dbPromise
  }

  async get<T>(key: string): Promise<T | undefined> {
    if (typeof window === 'undefined' || !('indexedDB' in window)) {
      throw new Error('indexedDB is not available')
    }

    const db = await this.openDb()
    return new Promise<T | undefined>((resolve, reject) => {
      const transaction = db.transaction(this.storeName)
      const objectStore = transaction.objectStore(this.storeName)
      const request = objectStore.get(key)
      request.onerror = () => reject(request.error)
      transaction.oncomplete = () => {
        resolve(request.result as T | undefined)
      }
    })
  }

  async set<T>(key: string, value: T): Promise<void> {
    if (typeof window === 'undefined' || !('indexedDB' in window)) {
      throw new Error('indexedDB is not available')
    }

    const db = await this.openDb()
    return new Promise<void>((resolve, reject) => {
      const transaction = db.transaction(this.storeName, 'readwrite')
      const objectStore = transaction.objectStore(this.storeName)
      objectStore.put(value, key)
      transaction.onerror = () => reject(transaction.error)
      transaction.oncomplete = () => {
        debugLog(`Stored ${key} locally`)
        resolve()
      }
    })
  }

  async del(key: string): Promise<void> {
    if (typeof window === 'undefined' || !('indexedDB' in window)) {
      throw new Error('indexedDB is not available')
    }

    const db = await this.openDb()
    return new Promise<void>((resolve, reject) => {
      const transaction = db.transaction(this.storeName, 'readwrite')
      const objectStore = transaction.objectStore(this.storeName)
      objectStore.delete(key)
      transaction.onerror = () => reject(transaction.error)
      transaction.oncomplete = () => {
        debugLog(`Deleted ${key} from local storage`)
        resolve()
      }
    })
  }
}

export default new GenericIndexedDbStore()
