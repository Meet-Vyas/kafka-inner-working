# Kafka — Learning How It Works Internally

[![progress-banner](https://backend.codecrafters.io/progress/kafka/7ffd0243-4924-460f-9a9f-4163ec119106)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This repository is my attempt to understand what happens behind the scenes when a client talks to Kafka. I am building a small Kafka-compatible server in Go as part of the [CodeCrafters "Build Your Own Kafka" challenge](https://codecrafters.io/challenges/kafka).

This is not a complete Kafka broker yet. The goal is to learn the basics step by step instead of using Kafka as a black box.

## What I have built so far

The server currently:

- Listens for TCP connections on port `9092`, which is Kafka's usual port.
- Accepts more than one client by handling each connection in its own goroutine.
- Reads the request size sent in the first four bytes of a Kafka message.
- Reads the API version and correlation ID from the request header.
- Sends the same correlation ID back so the client can match the response to its request.
- Returns error code `35` when the client asks for an unsupported API version.
- Responds to an `ApiVersions` request and says that API versions `0` through `4` are supported.
- Includes Kafka's throttle-time and tagged-fields values in the response.

## How it works in simple terms

A Kafka request is a group of bytes sent over a TCP connection. Those bytes need to be read in a specific order because each part has a meaning.

My server first reads four bytes to find out how large the incoming message is. It then reads the request header, takes out the API version and correlation ID, and builds a response in the byte format Kafka expects. All numbers are written using big-endian byte order, which means the most important byte comes first.

The correlation ID works a little like a ticket number. A client puts it in a request, and the server returns the same number in the response. This helps the client know which request the response belongs to.

## What I learned

While building this, I learned:

- How a TCP server listens for and accepts connections.
- How goroutines allow multiple connections to be handled without making every client wait for the previous one.
- Why reading the exact number of bytes matters when working with a network protocol.
- How binary data is converted between byte slices and Go integers.
- What big-endian encoding means in practice.
- How Kafka frames messages using a size at the beginning of each request and response.
- Why correlation IDs are useful when several requests may be in progress.
- How Kafka reports unsupported protocol versions with an error code.
- How compact arrays, tagged fields, and throttle time appear in a Kafka response.
- That even a small protocol response requires every byte to be placed in exactly the right position.

The biggest takeaway for me is that Kafka communication is not magic. Underneath the higher-level client libraries, it is a carefully structured exchange of bytes over a normal TCP connection.

## Run it locally

You need Go `1.24` or newer.

```sh
./your_program.sh
```

The server will start listening on `0.0.0.0:9092`.

## Current limitations

This is still a learning project. It currently handles one request per connection and only contains the `ApiVersions` work completed so far. Features such as producing messages, fetching records, topics, partitions, storage, replication, and consumer groups are not implemented yet.
