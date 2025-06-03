import type { AlbumMetadata, ArtistMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { getAlbumTracks, getArtistTracks } from '../composables/fetchFromBackend'
import { setCurrentlyPlayingTrack } from '../composables/globalState'
import { trackWithImageUrl } from '../composables/logic'

export async function play(artist?: ArtistMetadata, album?: AlbumMetadata, track?: TrackMetadata | TrackMetadataWithImageUrl) {
  if (track) {
    setCurrentlyPlayingTrack(trackWithImageUrl(track))
  }
  else if (album) {
    const tracks = await getAlbumTracks(album.musicbrainz_album_id)
    setCurrentlyPlayingTrack(tracks[0])
  }
  else if (artist) {
    const tracks = await getArtistTracks(artist.musicbrainz_artist_id)
    setCurrentlyPlayingTrack(tracks[0])
  }
}
