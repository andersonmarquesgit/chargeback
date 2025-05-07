#!/bin/bash

echo "⏳ Waiting for Cassandra to be ready..."
until cqlsh cassandra -e "DESCRIBE KEYSPACES"; do
  >&2 echo "Cassandra not ready yet - sleeping 5s"
  sleep 5
done

echo "Cassandra is up - applying schema..."
cqlsh cassandra -f /docker-entrypoint-initdb.d/schema.cql
