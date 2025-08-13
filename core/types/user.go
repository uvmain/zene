package types

type SubsonicUser struct {
	Folders             []int  `json:"folder" xml:"folder,attr"`                           // Optional: IDs of music folders the user can access.
	Username            string `json:"username" xml:"username,attr"`                       // Required: The name of the new user.
	Email               string `json:"email" xml:"email,attr"`                             // Required: The email address of the new user.
	AdminRole           bool   `json:"adminRole" xml:"adminRole,attr"`                     // Optional: Admin privileges. Default: false
	ScrobblingEnabled   bool   `json:"scrobblingEnabled" xml:"scrobblingEnabled,attr"`     // Optional: Enable scrobbling. Default: true
	StreamRole          bool   `json:"streamRole" xml:"streamRole,attr"`                   // Optional: Play files. Default: true
	SettingsRole        bool   `json:"settingsRole" xml:"settingsRole,attr"`               // Optional: Change personal settings/password. Default: true
	JukeboxRole         bool   `json:"jukeboxRole" xml:"jukeboxRole,attr"`                 // Optional: Play in jukebox mode. Default: false
	DownloadRole        bool   `json:"downloadRole" xml:"downloadRole,attr"`               // Optional: Download files. Default: false
	UploadRole          bool   `json:"uploadRole" xml:"uploadRole,attr"`                   // Optional: Upload files. Default: false
	PlaylistRole        bool   `json:"playlistRole" xml:"playlistRole,attr"`               // Optional: Create/delete playlists. Default: false
	CoverArtRole        bool   `json:"coverArtRole" xml:"coverArtRole,attr"`               // Optional: Change cover art/tags. Default: false
	CommentRole         bool   `json:"commentRole" xml:"commentRole,attr"`                 // Optional: Create/edit comments/ratings. Default: false
	PodcastRole         bool   `json:"podcastRole" xml:"podcastRole,attr"`                 // Optional: Manage podcasts. Default: false
	ShareRole           bool   `json:"shareRole" xml:"shareRole,attr"`                     // Optional: Share files. Default: false
	VideoConversionRole bool   `json:"videoConversionRole" xml:"videoConversionRole,attr"` // Optional: Start video conversions. Default: false
	MaxBitRate          int64  `json:"maxBitRate" xml:"maxBitRate,attr"`                   // Optional: Maximum bitrate for streaming. Default: 0 (no limit)
}

type SubsonicUsers struct {
	User []SubsonicUser `xml:"user" json:"user"`
}

type User struct {
	Id                  int64  `json:"id"`
	Username            string `json:"username" xml:"username"`                       // Required: The name of the new user.
	Password            string `json:"password" xml:"password"`                       // Required: The password of the new user, either clear text or hex-encoded.
	Email               string `json:"email" xml:"email"`                             // Required: The email address of the new user.
	LdapAuthenticated   bool   `json:"ldapAuthenticated" xml:"ldapAuthenticated"`     // Optional: LDAP authentication. Default: false
	AdminRole           bool   `json:"adminRole" xml:"adminRole"`                     // Optional: Admin privileges. Default: false
	ScrobblingEnabled   bool   `json:"scrobblingEnabled" xml:"scrobblingEnabled"`     // Optional: Enable scrobbling. Default: true
	SettingsRole        bool   `json:"settingsRole" xml:"settingsRole"`               // Optional: Change personal settings/password. Default: true
	StreamRole          bool   `json:"streamRole" xml:"streamRole"`                   // Optional: Play files. Default: true
	JukeboxRole         bool   `json:"jukeboxRole" xml:"jukeboxRole"`                 // Optional: Play in jukebox mode. Default: false
	DownloadRole        bool   `json:"downloadRole" xml:"downloadRole"`               // Optional: Download files. Default: false
	UploadRole          bool   `json:"uploadRole" xml:"uploadRole"`                   // Optional: Upload files. Default: false
	PlaylistRole        bool   `json:"playlistRole" xml:"playlistRole"`               // Optional: Create/delete playlists. Default: false
	CoverArtRole        bool   `json:"coverArtRole" xml:"coverArtRole"`               // Optional: Change cover art/tags. Default: false
	CommentRole         bool   `json:"commentRole" xml:"commentRole"`                 // Optional: Create/edit comments/ratings. Default: false
	PodcastRole         bool   `json:"podcastRole" xml:"podcastRole"`                 // Optional: Manage podcasts. Default: false
	ShareRole           bool   `json:"shareRole" xml:"shareRole"`                     // Optional: Share files. Default: false
	VideoConversionRole bool   `json:"videoConversionRole" xml:"videoConversionRole"` // Optional: Start video conversions. Default: false
	MaxBitRate          int64  `json:"maxBitRate" xml:"maxBitRate"`                   // Optional: Maximum bitrate for streaming. Default: 0 (no limit)
	Folders             []int  `json:"folder" xml:"folder"`                           // Optional: IDs of music folders the user can access.
}
