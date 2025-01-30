📌 Como Subir o Docker Compose
1️⃣ Construir e subir os containers

docker-compose up --build -d

2️⃣ Ver logs da aplicação

docker logs -f transaction-service

3️⃣ Conectar manualmente ao banco (caso queira testar)

docker exec -it postgres-db psql -U postgres -d dbname

docker exec -it postgres-db psql -U postgres -d transactions -c "\dt"

docker exec -it postgres-db psql -U postgres -d transactions -c "SELECT * FROM operation_types;"
