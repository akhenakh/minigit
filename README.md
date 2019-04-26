# minigit
A minimalist git clone (clone and pull) in pure Go
 
Main use is for the http-git Caddy plugin https://caddyserver.com/docs/http.git, so a pure Go image can be generated without cgo.

It supports github access token using the `-ghtoken` flag.
