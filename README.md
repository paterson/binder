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

## Setup

Install the project locally:
```bash
mkdir -p $GOPATH/src/github.com/paterson
cd $GOPATH/src/github.com/paterson
git clone github.com/paterson/binder
cd binder
```

Build and run the Project:
```bash
make
```

This sets up the authentication server, the directory service and the two initial file servers. To then run the client, open a new tab and run:

```bash
client
```

To run the default client. 
