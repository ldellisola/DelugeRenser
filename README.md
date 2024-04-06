# DelugeRenser
DelugeRenser, meaning "Deluge Cleaner" in Norwegian, is a simple tool to remove torrents stuck in seeding for ages until their seeding goal is achieved.
For example, if you have a seeding goal of 1.0, TorrentRenser will remove the torrent after one month even if that goal has not been achieved.

This tool is only compatible with **deluge v2**.

## Configuration
There are several parameters that can be configured using environment variables:

- `DELUGE_HOSTNAME`: The hostname of the deluge server. Default: `localhost`
- `DELUGE_PORT`: The port of the deluge server RPC API. Default: `58846`
- `DELUGE_USERNAME`: The username to authenticate with the deluge server. Default: `localclient`
- `DELUGE_PASSWORD`: The password to authenticate with the deluge server. **Mandatory**
- `KEEP_FOR`: The duration to keep seeding torrents. Default: `720h`
- `RUN_EVERY`: The interval to run the cleanup job. Default: `24h`
- `DRY_RUN`: If set to `true`, the tool will only log the torrents that would be removed. Default: `false`

The username and password can be found in the `auth` file in the deluge configuration directory. It has the following format:
```
<username>:<password>:<userLevel>
```
More info [here](https://dev.deluge-torrent.org/wiki/UserGuide/Authentication).

## Usage
You can run this tool using the following docker-compose file:
```yaml
version: '3.7'

services:
  torrentrenser:
    image: ghcr.io/ldellisola/deluge-renser:latest
    environment:
      DELUGE_HOSTNAME: deluge
      DELUGE_PASSWORD: password
      KEEP_FOR: 720h
      RUN_EVERY: 24h
      DRY_RUN: "false"
    restart: unless-stopped

  deluge:
    // Your deluge container
```
