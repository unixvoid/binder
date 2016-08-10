# Binder
This space is the future home of *binder*, an API for storing/retrieving
binaries.  Binder will have a REST API exposed for submitting, retrieving,
and tagging binaries with metadata.  Once data is submitted it can be
viewed with the help of nginx and autoindexing.

### TODO Endpoints
`/register`  : register the user and assign admin key  
`/rotate`    : rotate the admin key  
`/upload`    : upload content to the server  
`/remove`    : remove content to the server

### Testing
`curl -i --form sec=R5QgQwB7qKT3OlaM568vLbGQb --form filename=testraft --form file=@raft.pdf https://cryo.unixvoid.com/upload`
