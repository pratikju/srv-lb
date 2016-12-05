[![Build Status](https://travis-ci.org/pratikju/srv-lb.svg?branch=master)](https://travis-ci.org/pratikju/srv-lb)
[![GoDoc](http://godoc.org/github.com/pratikju/srv-lb?status.png)](http://godoc.org/github.com/pratikju/srv-lb/lb)


# SRV Record Load Balancer library for Go

_SRV-LB_ is a load balancer library designed for use with service discovery solutions
that expose an interface through DNS SRV records
(e.g. [consul](https://consul.io/) or [skyDNS](https://github.com/skynetservices/skydns))


The library selects a `SRV` record answer according to specified load balancer algorithm,
resolves its `A` record to an ip, and returns an `Address` structure:

	type Address struct {
		Address string
		Port    uint16
	}


To select a DNS server you can use the value from your system's `resolv.conf` (the default),
specify it explicitly when configuring the library,
or set it as an ENV variable (e.g. `SRVLB_HOST=127.0.0.1:8600` to connect to a local consul agent) at run time.


The library defaults to use a "Round Robin" algorithm, but you can specify another or build your own (see below).


## Example:
### Default Load Balancer

	srvName := "foo.service.fligl.io"
	cfg, err := lb.DefaultConfig()
	if err != nil {
		panic(err)
	}

	l := lb.New(cfg, srvName)

	address, err := l.Next()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", address.String())
	// Output: 0.1.2.3:8001

- Uses dns server configured in `/etc/resolv.conf`
- Uses round robin strategy


### or build a generic load balancer

	srvName := "foo.service.fligl.io"
	cfg, err := lb.DefaultConfig()
	if err != nil {
		panic(err)
	}

	l := lb.NewGeneric(cfg)

	address, err := l.Next(srvName)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("%s", address.String())
	// Output: 0.1.2.3:8001

### or configure explicitly

	srvName := "foo.service.fligl.io"
	cfg := &lb.Config{
		Dns:      dns.NewLookupLib("127.0.0.1:8600"),
		Strategy: random.RandomStrategy,
	}
	l := lb.New(cfg, srvName)

	address, err := l.Next()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("%s", address.String())
	// Output: 0.1.2.3:8001



## Development
tests are run against some fixture dns entries I set up on fligl.io (`dig foo.service.fligl.io SRV`).

	go get -u -t ./...
	go test ./...

	
## Build your own load balancing strategy

`srv-lb` leverages go's `init()` function to allow you to use your own
load balancer strategy without forking the core library. Below is a walkthrough
of how to create your own "FancyLB" strategy. For a complete example,
[see how the "random" strategy is implemented](https://github.com/pratikju/srv-lb/blob/master/strategy/random/random.go).

Of course, if your strategy would be generally usefull I would lolve a pull request!

_The default strategy, `RoundRobin`, is registered slightly differently to avoid import cycles, so avoid using it as an example_


Give your strategy a unique identifier

	const FancyStrategy lb.StrategyType = "fancy"

Create a factory (and your implementation of `GenericLoadBalancer`)

	func New(lib dns.Lookup) lb.GenericLoadBalancer {
		return &FancyLB{Dns: lib}
	}

Register it with the load balancer

	func init() {
		lb.RegisterStrategy(FancyStrategy, New)
	}


And then specify it when constructing your load balancer

	cfg, _ := lb.DefaultConfig()
	cfg.Strategy = fancy.FancyStrategy
	
	l := lb.New(cfg, srvName)


## Projects Using SRV-LB

- [chinchilla](https://github.com/pratikju/chinchilla) - a rabbitmq to REST bridge
