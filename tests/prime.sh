#!/bin/bash

number=$(shuf -i 0-12 -n1)

i=2

echo ${SCRIPT_TIMEOUT:-5s}
sleep ${SCRIPT_TIMEOUT:-5s}

if [ $number -lt 2 ]
then
    echo "$number IS NOT A PRIME NUMBER."
    exit
fi

while [ $i -lt $number ]
do
    if [ `expr $number % $i` -eq 0 ]
    then
        echo "$number IS NOT A PRIME NUMBER."
        exit 1
    fi
    i=`expr $i + 1`
done

echo "$number IS A PRIME NUMBER."