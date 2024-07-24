#!/usr/bin/env sh


set -e

host="$1"
shift
port="$1"
shift
cmd="$@"

until nc -z "$host" "$port"; do
  echo "Waiting for $host:$port to be ready..."
  sleep 1
done

>&2 echo "$host:$port is up - executing command"
exec $cmd
