# Zene
![Zene screenshot](./docs/assets/zene-home.webp)

## Self hosted Music Server and Web player
### Fast and feature packed with smart caching
- All transcoded audio is cached locally and cleaned with smart rules
- Fast full text search
- lyrics cached in database on retrieval
- Image endpoints and frontend make use of If-Modified-Since headers for 304 responses
- ffmpeg and ffprobe automatically downloaded as required on first boot
- album art automatically fetched from album folder || embedded in track || deezer || coverartarchive.org
- artist art automatically fetched from artist folder || deezer || wikidata

### Uses the OpenSubsonic API, with a few extras
Supports the following OpenSubsonic API extensions:
- `apiKeyAuthentication` (this project supports password, enc:password, salt & token, and ApiKey auth)
- `formPost` (all endpoints support GET and POST, with either formData values OR query parameters)
- `songLyrics` (lyrics are pulled from lrclib on request, and saved locally in the database for future calls)
- `transcodeOffset` (supports streaming from an offset)

### Supports the following OpenSubsonic API endpoints:

[Implemented OpenSubsonic API endpoints](./docs/implemented-opensubsonic-endpoints.md)

### additional custom API endpoints include:
- `createAvatar`
- `updateAvatar`
- `deleteAvatar`

The above endpoints enable dynamic functionality consistent with the existing getAvatar endpoint
- `getArtistArt`

### Tech stack
- `Sqlite` database
- `Go` backend
- `Vue` frontend (embedded during build)

- uses `Air` and `Vite` for HMR, and `Caddy` for SSL in local development

## localdev
### requirements
- Go v1.24+
- Node 22+

### install dependencies
First install npm dependencies (this will install the frontend workspace and the Caddy localdev utility)
```bash
npm i
```
Then install the Golang requirements
```bash
npm run setup
```

resolving Caddy cert issues on debian/ubuntu
- Ensure libnss3-tools is installed
  ```bash
  sudo apt install libnss3-tools
  ```
- Enable port-binding for the caddy binary
  ```bash
  sudo setcap CAP_NET_BIND_SERVICE=+eip node_modules/caddy-baron/caddy
  ```
- If you still get an ERR_CERT_AUTHORITY_INVALID error, run:
  ```bash
  certutil -d sql:$HOME/.pki/nssdb -A -t "C,," -n "Caddy Local Authority" -i ~/.local/share/caddy/pki/authorities/local/root.crt
  ```

## TODO
- [ ] add proper support for multiple music folders
- [x] use User.MaxBitRate to limit bitrate
- [ ] define an enum for allowed maxBitRate values to use in handlers
- [ ] use goose (or an alternative) to manage future database migrations
- [ ] getScanStatus and startScan handlers exist but need implementing
- [ ] HandleGetCoverArt and HandleGetArtistArt needs to handle size int param and resize if requested
