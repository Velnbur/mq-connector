# mq-connector

A small library for message brokers in DL projects.

## Key concepts

### Consumer And Producer

**Producer** - is an simple interface for publishing messages to queue.
Every **producer** is stiked to queue for which it was created, so you have
a **producer** for each proccessed queue.

**Consumer** - is long running proccess, that somehow handles every received
message from **producer**. 
 
```mermaid
flowchart LR
     idproducer1((Producer1)) --> idsendmessage1[Message1];
     idpruducer2((Producer2)) --> idsendmessage2[Message2];
     idsendmessage1[Message1] --> idqueue{Some Queue};
     idsendmessage2[Message2] --> idqueue{Some Queue};
     idqueue{Some Queue} --> idreceivedmessage1[Message1];
     idqueue{Some Queue} --> idreceivedmessage2[Message2];
     idreceivedmessage1[Message1] -->  idconsumer1((Consumer1));
     idreceivedmessage2[Message2] -->  idconsumer2((Consumer2));
```

In this situation each consumer receives one message, from queue. But
there is no such situation when the same message will be consumed two
or more times by different (or the same) **consumers**.

### Subscriber and Publisher

**Publisher** - is the the special case of **producer**, that publishes a message
not to a single queue, but to an _**subscription**_ (exchange or some kind of a router)
that will copy message to all subscribed (listening) queues.

**Subscriber** - is the special case of **consumer**, that receives messages from not
a specific queue but rather exchange or router that it subscribed to.
