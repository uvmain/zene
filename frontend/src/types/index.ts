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

export interface AlbumMetadata {
  artist: string
  album: string
  musicbrainz_track_id: string
  musicbrainz_album_id: string
  musicbrainz_artist_id: string
  genres: string[]
  release_date: string
  image_url: string
}
