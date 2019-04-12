rm ./schema.sql
mysqldump -u root messenger -p root > schema.sql
git add -A
git commit -m "update"
git push