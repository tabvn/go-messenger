#!/usr/bin/env bash
rm ./schema.sql
mysqldump -u root messenger -p  > schema.sql
git add -A
git commit -m "update"
git push
