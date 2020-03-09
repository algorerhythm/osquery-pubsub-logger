#!/bin/bash

sudo mkdir -p /usr/local/osquery_extensions
sudo cp etc/extensions.load /usr/local/osquery_extensions
sudo chown `whoami`:staff /usr/local/osquery_extensions/pubSubLogger.ext
go build -o /usr/local/osquery_extensions/pubSubLogger.ext
sudo chown -R root /usr/local/osquery_extensions/
