# Redis

Redis provider implementation.

Better encoder:

- Encode: `session.MSGPEncode`
- Decode: `session.MSGPDecode`

## Clients

The redis provider supports the upstream standard client via New(), the sentinel client via NewFailover(), and
the sentinel fail over client via NewFailoverCluster(). 

The difference between the standard client and the sentinel clients is that the standard client will only connect to a 
single redis server and the sentinel clients will connect to sentinel servers and figure out which redis server 
to connect to automatically. This allows for a failure of a redis server if you configure sentinel and redis correctly. 

The difference between the sentinel client via NewFailover() and the sentinel fail over client via NewFailoverCluster()
is the fail over client will fail over to other sentinels if they are configured in the event of a failure.