#!/bin/bash
set -ex

echo -n | sudo tee /var/lib/mysql/mysql-slow.log
echo -n | sudo tee /var/log/nginx/kataribe.log
