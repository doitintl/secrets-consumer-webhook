#!/bin/sh

echo "testing subtree each path as key name with a single value"
[ -n "$API_KEY" ] && echo "API_KEY: $API_KEY"
[ -n "$DATABASE_URL" ] && echo "DATABASE_URL: $DATABASE_URL"
[ -n "$DB_PASSWORD" ] && echo "DB_PASSWORD: $DB_PASSWORD"
[ -n "$APP_USER" ] && echo "APP_USER: $APP_USER"
[ -n "$DB_USER" ] && echo "DB_USER: $DB_USER"
exit 0
