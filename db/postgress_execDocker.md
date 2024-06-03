docker exec -it db-postgres-1 psql bank -U root

createdb --username=root --owner=root simple_bank
dropdb simple_bank

docker exec -it db-postgres-1 createdb --username=root --owner=root simple_bank
docker exec -it db-postgres-1 psql -U root simple_bank