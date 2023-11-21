#!/bin/bash
/usr/local/lsws/bin/lswsctrl start

while sleep 60; do
    if ! /usr/local/lsws/bin/lswsctrl status | grep 'litespeed is running' >/dev/null; then
        break
    fi
done

exit 1
