# grpcMessageBoard
Server and client for a simple message board. Unary and server-streaming RPCs.

## Server
- Stores messages with id, author name, title and body.
- When amount of messages exceedes `messageBoardServer.maxSize`, oldest messages are deleted.
- Safe concurrent operations with data using `sync.RWMutex`.
- A dockerfile is available for building the server's image.

## Client
- CLI app.
- Enter address of a running message board server to connect.
- Set your name.
- Post a message to the board.
- Show latest messages.