# minigit
A minimalist git clone (clone and pull) in pure Go
 
Main use is for the http-git Caddy plugin https://caddyserver.com/docs/http.git, so a pure Go image can be generated without cgo.

It supports github access token using the `-ghtoken` flag and several very limited command.

```
Usage:
  git [command]

Available Commands:
  clone       clone git repo
  help        Help about any command
  log         git log
  pull        pull git repo

Flags:
      --ghtoken string   github access token
```
