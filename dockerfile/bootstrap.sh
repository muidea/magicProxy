#!/bin/sh

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENERS ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -Address '$LISTENERS
fi

if [ $BROKER_LIST ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -Brokers '$BROKER_LIST
fi

if [ $METAAPI_SERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ApiSvr '$METAAPI_SERVER
fi

if [ $AUTHAPI_SERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -AuthSvr '$AUTHAPI_SERVER
fi

echo "Starting magicProxy..."

/var/app/magicProxy $EXTRA_ARGS "$@"