#!/bin/sh
/usr/local/bin/redis-server ./7000/redis.conf --pidfile ./7000/redis.pid --logfile ./7000/redis.log --cluster-config-file ./7000/node.conf&
/usr/local/bin/redis-server ./7001/redis.conf --pidfile ./7001/redis.pid --logfile ./7001/redis.log --cluster-config-file ./7001/node.conf&
/usr/local/bin/redis-server ./7002/redis.conf --pidfile ./7002/redis.pid --logfile ./7002/redis.log --cluster-config-file ./7002/node.conf&

./redis-trib.rb create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002
