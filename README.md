# CloudStorageSystem
Server and Client code of a cloud storage system

## Usage
1. Run the server:
```shell
go run cmd/SurfstoreServerExec/main.go -s <service> -p <port> -l -d (BlockStoreAddr*)
```
- `service`: The server side has 2 parts which are a Meta Server and a Block Store Server. Files are divided into blocks and stored in the Block Store Server. Meta Server contains informations where the files' blocks are and the versions of files used to synchronize the updates of multiple clients. Here, `service` should be one of three values: meta, block, or both. This is used to specify the service provided by the server. 
- `port`: This option defines the port number that the server listens to (default=8080). 
- `-l`: It configures the server to only listen on localhost. 
- `-d`: It configures the server to output log statements. 
- `(BlockStoreAddr\*)`: This is the BlockStore address that the server is configured with. If `service=both` then the BlockStoreAddr should be the `ip:port` of this server.

2. Run the client:
```shell
go run cmd/SurfstoreClientExec/main.go -d <meta_addr:port> <base_dir> <block_size>
- `meta_addr:port`: The IP address and port number of the Meta Server
- `base_dir`: This option indicates the directory to be synchronized with the cloud.
- `block_size`: This is the size of blocks stored on Block Store Server. Note that we make this size a mutable option in the command for the convenience of testing different sizes. But once you start your server, you should fix the block_size. If you want to change the block size, you need to restart the server.
