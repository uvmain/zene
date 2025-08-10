# Zene
![Zene screenshot](./docs/Screenshot%202025-08-10%20165519.png)

## Self hosted Music Server and Web player

### Uses the OpenSubsonic API, with a few extras
additional API endpoints include:
- createAvatar
- updateAvatar
- deleteAvatar
- getArtistArt


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

### resolve Caddy cert issues on debian/ubuntu
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
- [ ] use User.MaxBitRate to limit bitrate
- [ ] define an enum for allowed maxBitRate values to use in handlers
- [ ] use goose (or an alternative) to manage future database migrations
- [ ] getScanStatus and startScan handlers exist but need implementing
- [ ] HandleGetCoverArt and HandleGetArtistArt needs to handle size int param and resize if requested
