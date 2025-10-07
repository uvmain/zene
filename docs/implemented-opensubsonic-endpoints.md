## System
- [x] ping
- [x] getLicense
- [x] getOpenSubsonicExtensions
- [x] tokenInfo
## Browsing
- [x] getMusicFolders
- [x] getIndexes
- [x] getMusicDirectory
- [x] getGenres
- [x] getArtists
- [x] getArtist
- [x] getAlbum
- [x] getSong
- [x] getVideos[^1]
- [x] getVideoInfo[^1]
- [x] getArtistInfo[^2]
- [x] getArtistInfo2[^2]
- [x] getAlbumInfo[^4]
- [x] getAlbumInfo2[^4]
- [x] getSimilarSongs
- [x] getSimilarSongs2
- [x] getTopSongs[^5]
## Album/song lists
- [x] getAlbumList[^8][^9]
- [x] getAlbumList2[^8][^9]
- [x] getRandomSongs[^7][^9]
- [x] getSongsByGenre
- [x] getNowPlaying
- [x] getStarred
- [x] getStarred2
## Searching
- [x] search
- [x] search2
- [x] search3
## Playlists
- [x] getPlaylists
- [x] getPlaylist
- [x] createPlaylist
- [x] updatePlaylist[^6]
- [x] deletePlaylist
## Media retrieval
- [x] stream
- [x] download
- [ ] hls
- [x] getCaptions[^1]
- [x] getCoverArt
- [x] getLyrics
- [x] getAvatar
- [x] getLyricsBySongId
## Media annotation
- [x] star
- [x] unstar
- [x] setRating
- [x] scrobble[^3]
## Sharing
- [ ] getShares
- [ ] createShare
- [ ] updateShare
- [ ] deleteShare
## Podcast
- [x] getPodcasts
- [x] getPodcastEpisode
- [x] getNewestPodcasts
- [x] refreshPodcasts
- [x] createPodcastChannel
- [x] deletePodcastChannel
- [x] deletePodcastEpisode
- [x] downloadPodcastEpisode
## Jukebox
- [x] jukeboxControl[^1]
## Internet radio
- [x] getInternetRadioStations
- [x] createInternetRadioStation
- [x] updateInternetRadioStation
- [x] deleteInternetRadioStation
## Chat
- [x] getChatMessages
- [x] addChatMessage
## User management
- [x] getUser
- [x] getUsers
- [x] createUser
- [x] updateUser
- [x] deleteUser
- [x] changePassword
## Bookmarks
- [x] getBookmarks
- [x] createBookmark
- [x] deleteBookmark
- [x] getPlayQueue
- [x] getPlayQueueByIndex
- [x] savePlayQueue
- [x] savePlayQueueByIndex
## Media library scanning
- [x] getScanStatus
- [x] startScan

[^1]: Endpoint exists but returns an unsupported error.
[^2]: Similar artists are fetched from Deezer, not lastfm. Biography is not supported.
[^3]: Scrobble updates local Now Playing and Play Count - it does not integrate with lastfm.
[^4]: Notes property is not supported.
[^5]: Top songs are fetched from Deezer, not lastfm.
[^6]: Additionally allows `coverArt` and multiple `allowedUserId` params to be sent.
[^7]: Additionally support `offset` param to enable paging through the same random results.
[^8]: Additionally supports a `type` param value of `release`, ordering by release date desc.
[^9]: Additionally supports a `seed` integer param value for deterministic random ordering.