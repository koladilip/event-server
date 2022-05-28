# event-server
Sample Event Delivery system using GoLang based on these [requirements](https://github.com/koladilip/event-server/blob/main/EventDelivery.pdf).

## Architecture Diagram
![EventDeliverySystem](https://github.com/koladilip/event-server/blob/main/EventDeliverySystem.png)

## Components]
### API 
Receives the events from various customer portals and stores then in Kafka for futher processing. Source events can be stored in several partitions but we need to store all events of a user ID in the same partition else we can't maintain order.

### Source Events Reader
Reads events stored in Kafka topic and hand overs them to transformer.

### Transformer
Transforms the source events to supported destination event types and stores them in destination topics (every destination has a dedicated topic)

### Destinations
We need to support several destinations and each destination will have its own format and protocals for receiving events. We have defined high level interface for destination to abstract out the custom implementation for each destination.
```go
type Destination interface {
	Id() string
	Transform(event.SourceEvent) (event.DestinationEvent, error)
	Deliver(event.DestinationEvent) error
}
```
In order to support Delivery isolation, we are storing destination events in separate topics(one topic will be dedicated for one destination). This way if destination is down then it won't effect other destinations and also we need to maintain order of user events so we make sure that all events of a user ID goes to same parition within the destination topic.

### Deliverers
Delivers events from each destination topic to corresponding endpoint.

## How to run locally?
1. Install and start [kafka](https://kafka.apache.org/quickstart)
1. `git clone https://github.com/koladilip/event-server.git`
1. `cd event-server`
1. Run server: `rm -rf server.log && go run main.go > server.log 2>&1`
1. Run test script: `go run test/main.go`
1. Check results: `tail -f server.go`