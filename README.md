# event-server
Sample Event Delivery system using GoLang based on these [requirements](https://github.com/koladilip/event-server/blob/main/EventDelivery.pdf).

## Architecture Diagram
![EventDeliverySystem](https://github.com/koladilip/event-server/blob/main/EventDeliverySystem.png)

## Components]
### API 
Receives the events from various customer portals and stores then in Kafka for futher processing.

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

### Deliverers
Delivers event from each destination topic to the destination endpoint.

## How to run locally?
1. Install and start [kafka](https://kafka.apache.org/quickstart)
1. `git clone https://github.com/koladilip/event-server.git`
1. `cd event-server`
1. Run server: `rm -rf server.log && go run main.go > server.log 2>&1`
1. Run test script: `go run test/main.go`
1. Check results: `tail -f server.go`