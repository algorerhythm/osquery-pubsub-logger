#!/bin/bash
sudo GOOGLE_APPLICATION_CREDENTIALS=/var/osquery/certs/osquery-creds.json SOCKET_PATH=/var/osquery/osquery.em TOPIC=osquery GCP_PROJECT=notebook-269803 osqueryd --flagfile=etc/osquery.flags
