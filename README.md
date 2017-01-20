# CS4032: Project

### Distributed File System

## Milestones

- [x] Token Based Authentication Service
- [x] Distributed File Servers
- [x] Directory Service to co-ordinate the file system
- [x] Encryption of all data transferred
- [x] Replication
- [x] File locking to ensure correctness
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

### Token Based Authentication System
I used a 3 key security model that can be used to mutually authenticate clients and servers. This ensures that all communication between all system components are encrypted. Once the client logs in using a username and password, the authentication system responds with a token comprised of a ticket, a session key, the identity of the server the ticket is for, and a timeout period for the ticket. The ticket contains a copy of the session key. The session key is a key generated at random (ie. when the client requests to login) to be used to encrypt and decrypt communication between the client and server. The ticket is itself encrypted with a server encryption key, which is known only by the authentication system and the server the client wants to contact. Thus only the server can decrypt the ticket. The token is encrypted with a key derived from the client’s password. Thus, the token can only be decrypted by the client (assuming the password is secure).


To send requests to the server the client now uses the session key to encrypt all messages before sending them to the server. It sends the ticket (without encrypting it) to the server along with each request (i.e. request = message encrypted + ticket). The server, on receiving the request, decrypts the ticket to obtain the session key. It then decrypts the message with that session key. It performs the required action, assuming successful decryption, generates a response and encrypts the response with the session key. The client can decrypt the response as it knows the session key.

This is an effective and secure system for a distributed system. Once the system components all have a shared server key, they can each individually authenticate a request.

### Distributed File Servers
The system can utilise an unlimited number of file servers, and adding a new file server is as simple as adding it to the directory service's database. Any number of file servers can be spun up quickly using docker.  

### Replication
Any files added to the system will automatically be replicated to a number of other file servers (up to 4) so that the system will continue to function if a file server goes down. 

### Directory Service to Co-ordinate The File System
The directory service co-ordinates the file system - it keeps track of the the file servers, routes read and write requests, handles locking and manages replication of file servers. The directory service also ensures that the system remains balanced - the files are evenly distributed over the file servers by allocating new files to the least used file servers. 

### File Locking

The system allows Write Locking of files, which ensures correctness of the system. Reading files do not require a lock. 

### Client Caching

When a file is read to a client, it is added to the client's 2MB local cache, which automatically evicts entries after a timeout, or when the cache is full (it evicts the oldest entry then) 

### Client Proxy

The Client Proxy is a library to abstract interaction with the system, it allows you login/signup, request read and write permissions, read, write, lock and unlock files.
