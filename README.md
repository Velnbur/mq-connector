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

```mermaid
flowchart LR
    idpublisher1((Publisher1)) --> idsendmessage1[Message1];
    idpublisher2((Publisher2)) --> idsendmessage2[Message2];
    idsendmessage1[Message1] --> idsubs[(Subscription)];
    idsendmessage2[Message2] --> idsubs[(Subscription)];
    idsubs[(Subscription)] --> idreceivedmessage11[Message1];
    idsubs[(Subscription)] --> idreceivedmessage12[Message1];
    
    subgraph queue1 [some queue1]
        idreceivedmessage11[Message1] --- idreceivedmessage21[Message2];
    end
    
    subgraph queue2 [some queue2]
        idreceivedmessage12[Message1] --- idreceivedmessage22[Message2];
    end
    
    idreceivedmessage21[Message2] --> idsubscriber1((Subscriber1));
    idreceivedmessage22[Message2] --> idsubscriber2((Subscriber2));
```

### Request/Reply model `RpcServer`, `RpcClient`

Here, for communication between two services, two queues will be created
for requests and responses respectivly. Each request message will have a
`ReplyTo` parameter, where will be the name of the queue to which send the
reponse message.

```mermaid
flowchart LR
    idservice1((Service1)) --> idrequest1[Request1];
    
    subgraph requests [Requests Queue]
        idrequest1[Request1] --- idrequest2[Request2];
        idrequest2[Request2] --- idrequest3[Request3];
    end
    
    idrequest3[Request3] --> idservice2((Service2));
    
    idservice2((Service2)) --> idresponse1[Response1];
    
    subgraph responses [Responses Queue]
        idresponse1[Response1] --- idresponse2[Response2];
        idresponse2[Response2] --- idresponse3[Response3];
    end
    
    idresponse3[Response3] --> idservice1((Service1));
```
