import type { SubsonicSong } from '~/types/subsonicSong'

export const routeTracks = ref<SubsonicSong[]>([])

export function clearRouteTracks() {
  routeTracks.value = [] as SubsonicSong[]
}
