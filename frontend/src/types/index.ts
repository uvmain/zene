export interface SlugWithDimensions {
  slug: string
  width: number
  height: number
}

export interface ImageMetadata {
  filePath: string
  fileName: string
  title: string
  dateTaken: string
  dateUploaded: string
  cameraMake: string
  cameraModel: string
  lensMake: string
  lensModel: string
  fStop: string
  exposureTime: string
  flashStatus: string
  focalLength: string
  iso: string
  exposureMode: string
  whiteBalance: string
  whiteBalanceMode: string
  width: number
  height: number
  orientation: string
  panoramic: boolean
}

export interface GenreMetadata {
  genre: string
  count: string
}

export interface ArtistMetadata {
  musicbrainz_artist_id: string
  artist: string
  image_url: string
}

export interface AlbumMetadata {
  artist: string
  album_artist: string
  album: string
  musicbrainz_album_id: string
  musicbrainz_artist_id: string
  genres: string[]
  release_date: string
  image_url: string
}

export interface TrackMetadata {
  file_path: string
  file_name: string
  date_added: string
  date_modified: string
  format: string
  duration: string
  size: string
  bitrate: string
  title: string
  artist: string
  album: string
  album_artist: string
  genre: string
  track_number: string
  total_tracks: string
  disc_number: string
  total_discs: string
  release_date: string
  musicbrainz_artist_id: string
  musicbrainz_album_id: string
  musicbrainz_track_id: string
  label: string
  user_play_count: number
  global_play_count: number
}

export interface TrackMetadataWithImageUrl extends TrackMetadata {
  image_url: string
}

export interface Queue {
  tracks: TrackMetadataWithImageUrl[]
  position: number
}

export interface StandardResponse {
  status: string
  error?: string
}
