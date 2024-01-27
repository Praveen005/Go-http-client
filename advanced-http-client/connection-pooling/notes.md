# Connection Pooling

In header-middleware & logging-middleware(see the respective folders), we saw the default RoundTripper
interface implementation carries an HTTP request to the remote
server and then carries the response back.

One of the underlying steps that occurs is that a new TCP connection
is established for your request. This connection setup process is
expensive. You may not notice it when you are making a single
request. However, when you are making HTTP requests as part of a
service-oriented architecture, for example, you are usually making
multiple requests in a short time window—either in bursts or
continuously. In such a scenario, it is expensive to perform the TCP
connection setup for every request. Hence, the net/http library
maintains a connection pool where it automatically tries to reuse an
existing TCP connection to send your HTTP requests.

The `net/http/httptrace` package will help us delve into the internals of connection pooling.



## DNS Resolution

DNS resolution, or Domain Name System resolution, is the process of translating a human-readable domain name into an IP address that computers can understand. When you type a URL into your web browser, the browser needs to know the IP address of the server hosting the website. This is where DNS resolution comes in.
Here's a simplified explanation of the process:

1. **User request:** You type a URL into your web browser, such as `www.example.com`.
2. **DNS query:** Your computer sends a request to a DNS resolver, asking for the IP address associated with the domain name `www.example.com`.
3. **DNS response:** The DNS resolver looks up the IP address in its cache or queries other DNS servers until it finds the IP address. It then sends the IP address back to your computer.
4. **Browser connects:** Your web browser uses the IP address to connect to the server hosting the website, and the website's content is loaded in your browser.

## httptrace.GotConn

In the context of the `httptrace.GotConn` function in Go, it is specifically triggered when the TCP connection is successfully established, providing information about the connection, such as the remote and local addresses, whether the connection was reused, and the time of establishment. 


`Connection establishment` refers to the process of creating a communication link (connection) between the client (e.g., our application making an HTTP request) and the server (the destination of the HTTP request). The establishment of a connection involves several steps, including:

1. DNS Resolution: If the server's domain name is not already resolved to an IP address, a Domain Name System (DNS) resolution is performed to obtain the IP address associated with the server's domain.

2. TCP Handshake: Once the IP address is known, a Transmission Control Protocol (TCP) connection is initiated between the client and the server. This involves a three-way handshake:

    1. SYN (Synchronize): The client sends a SYN packet to the server, indicating its intention to establish a connection.

    2. SYN-ACK (Synchronize-Acknowledge): The server responds with a SYN-ACK packet, acknowledging the client's request and indicating its own readiness to establish the connection.

    3. ACK (Acknowledge): Finally, the client sends an ACK packet, acknowledging the server's response. At this point, the TCP connection is considered established.

3. TLS/SSL Handshake (if using HTTPS): If the communication is secure and is done over HTTPS, a Transport Layer Security (TLS) or Secure Sockets Layer (SSL) handshake may take place to establish a secure encrypted connection.

The HTTP request comes after the establishment of the TCP connection and, if applicable, the TLS/SSL handshake in the case of HTTPS. The purpose of the handshake is to set up the parameters of the TCP connection and, if needed, negotiate security features for the subsequent data exchange.

### TLS/SSL handshake

 In TLS (Transport Layer Security) or SSL (Secure Sockets Layer) handshake, there is an additional set of handshakes beyond the initial three-way TCP handshake. The TLS handshake involves a four-step process:

1. ClientHello: The client sends a "ClientHello" message to the server, indicating its intention to establish a secure connection and specifying the cryptographic algorithms and other parameters it supports.

2. ServerHello: The server responds with a "ServerHello" message, selecting the cryptographic algorithms and other parameters from the client's list that it also supports. Additionally, the server sends its digital certificate (if required for authentication) and other information needed for the key exchange.

3. Key Exchange: The client and server perform a key exchange to establish shared secret information used for encrypting and decrypting data during the secure communication. The key exchange method depends on the chosen cryptographic algorithms and the server's configuration.

4. Finished: Both the client and server exchange "Finished" messages to confirm that the handshake is complete. At this point, they have agreed upon a set of parameters for secure communication, and the encrypted session begins.