# pricer

Depends on generator https://github.com/chucky-1/generator

The pricer gets prices from the generator by redis stream.

Then he pushes all prices to the broker by grpc stream.

The client can subscribe to the prices he is interested in. Subscription implemented for GRPS.

The broker https://github.com/chucky-1/broker

The client https://github.com/chucky-1/trader