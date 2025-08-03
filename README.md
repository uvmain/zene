# localdev

## Caddy requirements on debian
- Ensure libnss3-tools is installed
  ```bash
  sudo apt install libnss3-tools
  ```
- Enable port-binding for the caddy binary
  ```bash
  sudo setcap CAP_NET_BIND_SERVICE=+eip node_modules/caddy-baron/caddy
  ```
