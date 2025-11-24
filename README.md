# broadcast-server

**broadcast-server** is a command-line WebSocket broadcast tool
 It allows multiple clients to connect to a server, and whenever any client sends a message, the server broadcasts it to all other connected clients in real time

This project provides:

-  Real-time messaging via WebSocket
-  Clean, user-friendly Cobra CLI
-  Graceful shutdown of server and clients
- Concurrent read/write loops to enable non-blocking message handling via goroutines

##  Usage

```
./broadcast-server start
```

WebSocket endpoint is `ws://localhost:9000/ws/connect`

```
./broadcast-server connect
```

## About This Project

This project was built to complete the **Broadcast Server** assignment on Roadmap.sh:

https://roadmap.sh/projects/broadcast-server

Developed by **qs-lzh**

## Support the Project

If you find this project useful, feel free to support it with an upvote on the Roadmap project page:

👉 https://roadmap.sh/projects/broadcast-server/solutions?u=6919ca0806aadfe789824b5c