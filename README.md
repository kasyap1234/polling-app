# polling-app
Overview of the Polling App (Voting System) Process
The code provided is for a real-time polling app where users can vote, and the votes are broadcasted to all connected clients in real-time using WebSockets. The application uses goroutines and channels for concurrent processing, ensuring efficient handling of multiple clients and real-time updates.

Here’s a high-level walkthrough of the entire process:

1. WebSocket Connection Establishment:
When a client (e.g., a browser) connects to the server, it establishes a WebSocket connection. This is done using the websocket.Upgrade method in the main function.
The server listens for incoming WebSocket connections at the /ws endpoint.
2. Hub Initialization:
The Hub is a central structure that manages the WebSocket connections and broadcasts messages to all connected clients.
The Hub has:
A map (clients) that holds all active clients.
A broadcast channel to send messages to all clients.
A register and unregister function to add or remove clients from the Hub.
3. Client Structure:
A Client represents an individual WebSocket connection. It contains:
conn: A WebSocket connection object.
send: A channel used to send messages to the client.
hub: A reference to the Hub, which manages all clients.
readPump() and writePump() are two methods that handle the reception and sending of messages for each client, respectively.
4. Client Registration:
When a client successfully connects, it is registered with the Hub via hub.register(client).
The Hub adds the client to its clients map and starts a goroutine to handle sending messages to that client (writePump).
5. Message Handling:
Once the connection is established, the client can send messages (votes) to the server.
The readPump function listens for incoming messages from the WebSocket connection. When a message is received, it:
Passes the message to the Hub via the hub.broadcast channel.
The broadcast channel is responsible for distributing the message to all connected clients.
6. Vote Counting:
The message received in the readPump function could be a vote (e.g., "Vote for Option A").
The Hub processes the vote and stores it in memory (not shown in the provided code, but typically a simple counter or a map).
Concurrency with Goroutines:
Since there can be multiple clients voting concurrently, goroutines handle each client independently without blocking others.
This ensures that the server can handle voting in parallel, allowing votes to be counted and broadcasted in real-time.
As each vote is processed, the Hub updates the vote count and broadcasts the updated vote count to all connected clients via the broadcast channel.
7. Broadcasting Updates:
The broadcast channel is used to send updates (e.g., new vote counts) to all connected clients.
The writePump function of each client listens for messages on its send channel and sends them to the WebSocket connection when received.
The writePump goroutine ensures that the client is continuously updated with the latest voting information as it is received from the Hub.
8. Handling Client Disconnections:
If a client disconnects or closes the WebSocket connection, the Hub removes the client from its clients map and closes the send channel associated with that client.
This is done in the unregister() function.
9. Concurrency and Synchronization:
Go routines are used to handle each client's reading (readPump) and writing (writePump) concurrently.
Channels are used to pass messages between the clients and the Hub in a thread-safe manner.
The Hub ensures that only one goroutine handles broadcasting to a client by using a separate send channel for each client.
10. Real-Time Updates:
The server continuously updates all connected clients with the latest voting data (vote counts, for example). As a client sends a vote, the vote is processed and immediately broadcasted to all clients.
Clients receive this data in real-time through their WebSocket connections.
In Summary: Voting Flow
A client connects to the server via WebSocket.
The Hub registers the client and starts listening for messages from it.
The client sends a vote (message) via the WebSocket connection.
The Hub processes the vote and updates the count.
The Hub broadcasts the updated vote count to all connected clients using the broadcast channel.
Clients receive the updates in real-time via their WebSocket connections.
The client disconnects when they close their WebSocket connection, at which point the Hub unregisters the client and cleans up resources.
This entire system ensures real-time, concurrent voting where votes are processed and broadcasted in real-time, allowing all clients to stay updated with the latest vote counts without needing to refresh their browsers.