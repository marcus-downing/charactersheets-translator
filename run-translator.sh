#!/bin/bash

if pgrep translator &> /dev/null; then
        killall translator
fi

cd /data/charsheets/charactersheets-translator
./translator &> /var/log/translator &
