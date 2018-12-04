rm ./schema.sql
mysqldump -u root messenger > schema.sql
git add -A
git commit -m "update"
git push