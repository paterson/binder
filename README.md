# CS4032: Project

### Distributed File System

## Milestones

- [x] Token Based Authentication Service
- [x] Distributed File Servers
- [x] Directory Service to co-ordinate the file system
- [x] Encryption of all data transferred
- [x] Replication
- [x] Client Caching
- [x] Client Proxy to abstract out communication with system
- [x] Client with example interactions

## Setup Locally

Install the project locally:
```bash
mkdir -p $GOPATH/src/github.com/paterson
cd $GOPATH/src/github.com/paterson
git clone github.com/paterson/binder
cd binder
```

Build and run the project:
```bash
make
```

This sets up the authentication server, the directory service and the two initial file servers. To then run the client, **open a new tab** and run:

```bash
client
```

To run the default client. You could also use `go get` to get each individual service, but this approach is simpler and encapsulates the nessecary steps.

## Setup on Docker

Due to the restrictions on SCSS OpenNebula, I can't create more than two VM's, which inhibits this, however each microservice is on Docker Hub and easily used:

```bash
go get paterson/binder-authservice
go get paterson/binder-directoryservice
go get paterson/binder-fileserver
```

If SCSS OpenNebula was changed, or if we were using this in real life, Docker would be an ideal way to deploy this system.

## Specifications

### Authentication Service
