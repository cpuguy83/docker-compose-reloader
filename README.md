# docker-compose-reloader
Monitor fs events and reload docker-compose on changes

Watches for FS events and reloads the docker-compose.yml from the current working dir.
Also sends livereload notifications.

### Usage

docker-compose-watcher --watch . --service app

if `--service` is not set reload all services.
If `--watch` is not set, watch current dir (and all subdirs).
`--watch` can be specified multiple times
