# gossip-gloomers

A series of [distributed programming riddles](https://fly.io/blog/gossip-glomers/) by fly.io

## Install dependencies

```sh
brew install openjdk graphviz gnuplot
```
Setup instructions can also be found [here](https://fly.io/dist-sys/1/) 

## Run Test

### Challenge #1: Echo

```sh
./run.sh 1_echo echo
```

### Challenge #2: Unique ID Generation

```sh
./run.sh 2_unique_ids unique-ids 3 --availability total --nemesis partition 
```

### Challenge #3: Unique ID Generation

#### Single Node
```sh
./run.sh 3_broadcast broadcast 1 10
```

#### Multi Node
```sh
./run.sh 3_broadcast broadcast 5 10
```

