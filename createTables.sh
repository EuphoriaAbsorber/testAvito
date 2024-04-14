printf "Copy scripts:\n"
sudo docker cp ./_postgres postgres:/
printf "DROP TABLE:\n"
sudo docker exec -it postgres psql -U postgres -d postgres -a -f ./_postgres/drop.sql -S | grep -E "(NOTICE|ERROR)" && printf "\n"
printf "CREATE TABLE:\n"
sudo docker exec -it postgres psql -U postgres -d postgres -a -f ./_postgres/create.sql | grep -E "(NOTICE|ERROR)" && printf "\n"