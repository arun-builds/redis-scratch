package config

var Host string = "0.0.0.0"
var Port int = 7379

var KeysLimit int = 100

// will evict EvictionRatio of keys whenever eviction runs
var EvictionRatio float64 = 0.40

var EvictionStrategy string = "allkeys-random"
var AOFFile string = "./redis-scratch.aof"
