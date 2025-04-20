# targeting_engine

# Run docker command

# for postgres 
 docker run -d --name postgres-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=906 -e POSTGRES_DB=postgres -p 5432:5432 -v pg_data:/var/lib/postgresql/data postgres:latest


## for Redis
docker run -d --name redis-server -p 6379:6379 -v redis_data:/data redis:latest redis-server --requirepass "906" --appendonly yes

## For GUI if required install
# for postgres use Dbeaver
# for redis use RedisInsight
# you can use any other of your preferred prefernce

# For metrics we are using prometheus anf grafana run 
docker compose up -d
# you can accces these
Prometheus: http://localhost:9090
Grafana: http://localhost:3000

## in grafana add the dashboard.json file which is present in metrics and add prometheus endpoint 