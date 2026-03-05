#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
  CREATE DATABASE app_prod;
  CREATE DATABASE app_test;
EOSQL   