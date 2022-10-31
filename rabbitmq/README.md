# rabbitmq

This subpackage is a **RabbitMQ** implementation of `mq-connector`'s interfaces.

## Usage

### In DL projects

Add config parser to `config` package by, firstly, adding necessary methods to
`Config` interface (see all available interfaces in `rabbitmq/config.go`).

```go
type Config interface {
    // ...
    
    
    rabbitmq.RabbitConsumerer
}
```

Secondly, add `RabbitConfiger` for config parsing to `config` struct:

``` go
type config struct {
    // ...
    *rabbitmq.RabbitConfiger
}
```

Thirdly, update `Config` constructor:

``` go
func New(getter kv.Getter) Config {
	return &config{
        // ...
        RabbitConfiger: rabbitmq.NewRabbitConfiger(getter),
    }
}
```

Fourthly, add rabbit connection configuration to `config.yaml`:

``` yaml
rabbit:
  url: "amqp://guest:guest@localhost:5672"
```

and add configuration of the needed object:

``` yaml
mb_consumer_yaml_key: # just an example. Key may be arbitary
  queue: "some_queue_name"
```

And know you may get `mqconnector.Consumer` from `cfg config.Config` with:

``` go
cfg.RabbitConsumer("mb_consumer_yaml_key") // just an example. Key may be arbitary 
```

