#!/bin/sh
set -e

# Generate ACL file from environment variables at startup.
# Disables the default user and creates a named user with full permissions.
{
  echo "user default off"
  echo "user ${REDIS_USERNAME} on >${REDIS_PASSWORD} ~* &* +@all"
} > /tmp/users.acl

exec redis-server --appendonly yes --aclfile /tmp/users.acl
