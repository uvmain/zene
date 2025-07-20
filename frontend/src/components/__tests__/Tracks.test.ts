import type { TrackMetadataWithImageUrl } from '../../types'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import Tracks from '../Tracks.vue'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

const mockTracks: TrackMetadataWithImageUrl[] = [
  {
    file_path: 'test/file/path.mp3',
    file_name: 'file.mp3',
    date_added: '2023-01-01',
    date_modified: '2023-01-01',
    format: 'mp3',
    duration: '180',
    size: '4 MB',
    bitrate: '320 kbps',
    title: 'Test Track',
    artist: 'Test Artist',
    album: 'Test Album',
    album_artist: 'Test Album Artist',
    genre: 'Test Genre',
    track_number: '1',
    total_tracks: '10',
    disc_number: '1',
    total_discs: '1',
    release_date: '2023-01-01',
    musicbrainz_artist_id: 'test-artist-id',
    musicbrainz_album_id: 'test-album-id',
    musicbrainz_track_id: 'test-track-id',
    label: 'Test Label',
    user_play_count: 0,
    global_play_count: 0,
    image_url: 'https://example.com/image.jpg',
  },
]

describe('tracks', () => {
  it('should render correctly', () => {
    const wrapper = mount(Tracks, {
      props: { tracks: mockTracks },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          RouterLink: true,
          RouterView: true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(Tracks, {
      props: { tracks: mockTracks },
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          RouterLink: true,
          RouterView: true,
        },
      },
    })
    expect(wrapper.vm).toBeTruthy()
  })
})
