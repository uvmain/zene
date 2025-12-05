import * as localforage from 'localforage'

export function usePodcastStore() {
  const setStoredEpisode = async (key: string, episode: Blob): Promise<void> => {
    try {
      await localforage.setItem(key, episode)
    }
    catch (error) {
      const err = error as Error
      console.error(`Error storing podcast episode in IndexedDB: ${err.message}`)
    }
  }

  const getStoredEpisode = async (key: string): Promise<Blob | null> => {
    try {
      const value = await localforage.getItem(key)
      return value as Blob | null
    }
    catch (error) {
      const err = error as Error
      console.log(err.message)
      return null
    }
  }

  const deleteStoredEpisode = async (key: string): Promise<void> => {
    try {
      await localforage.removeItem(key)
    }
    catch (error) {
      const err = error as Error
      console.error(`Error deleting podcast episode from IndexedDB: ${err.message}`)
    }
  }

  return {
    setStoredEpisode,
    getStoredEpisode,
    deleteStoredEpisode,
  }
}
