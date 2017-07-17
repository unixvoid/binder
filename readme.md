# Binder
[![Build Status (Travis)](https://travis-ci.org/unixvoid/binder.svg?branch=master)](https://travis-ci.org/unixvoid/binder)  
Binder is a project for the storage/retrieval of files.  As the name *bin*dir
implies, its main use in my concern is the storage of binaries.  I use travis-ci
as my main source of continous integration an I wanted to capture the binaries
that travis compiles for me.  Every time I make a git commit it triggers travis
to build the project and travis in turn pushes the statically compiled binaries
to binder.  You can see the project in work over
[here](https://cryo.unixvoid.com).  
  
Another use of this project is to auto push ACI/rkt images for public
consumption.  Whenever someone does a `rkt fetch unixvoid.com/<project>`, rkt
will be redirected to [binder](https://cryo.unixvoid.com/bin/rkt) for the aci
fetch.  
  
The project is intended ony for public storage as any file that is uploaded is
pushed publicly to the nginx-backed ui.  When binder starts up for the first
time it generates a security key that is needed for uploads.  You can only
upload/remove files from binder if you posess the security key.  
  
Yet another feature is the private key storage.  Private key storage is the storage
of keys that are not pubicly accessable, and require the security token to
upload/remove/retrieve.  I use this mainly for for storing private keys.  When I
have travis auto-build my rkt containers it needs to GPG sign the image and
upload it and the signed public key.  To accomplish these things, travis needs
my cooresponding GPG private key.  It does this by fetching it from travis.
This feature along with the public storage will be explained in more detail
below.


## Running binder
There are 3 main ways to run binder:

1. **From Source**: To run binder from source we need to pull the required
   golang dependencies and then run.  We cann accomplish this with the following
   make commands:  
   `make dependencies` `make run`  
   To compile the project statically, we use:  
   `make dependencies` `make stat`  
   This will produce a statically compiled binary in the `bin/` directory.

2. **Docker**: We have binder pre-packaged as a docker container over on the
   [dockerhub](https://hub.docker.com/r/unixvoid/binder), go grab the latest and
   run with: `docker run -d -p 8000:8000 unixvoid/binder`.  The binder can also
   be run in a docker-compose stack by executing `make compose`.

3. **ACI/rkt**: WIP


## Documentation
All documentation is in the [wiki](https://unixvoid.github.io/binder_docs/)
* [API](https://unixvoid.github.io/binder_docs/api/)
* [Configuration](https://unixvoid.github.io/binder_docs/configuration/)
