#!/bin/bash
set -ex

> /var/lib/mysql/mysql-slow.log
echo -n | sudo tee /var/log/nginx/kataribe.log
