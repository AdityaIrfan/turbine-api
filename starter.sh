#!/bin/bash

echo
echo $"NEEDS 1 ARGUMENT, './starter.sh develop' OR './starter.sh production'"
echo

# Arguments
server=$1

docker_compose_yml=""
env=""
# Set variable based on server
if [ "$server" = "develop" ]; then
  docker_compose_yml="docker-compose-develop.yml"
  env=".env.develop"
elif [ "$server" = "production" ]; then
  docker_compose_yml="docker-compose-develop.yml"
  env=".env.production"
else
  echo $'\xE2\x9C\x98' " UNAVAILABLE SERVER TYPE"
  exit 1
fi

# Function to check if a container exists
container_exists() {
  local container_name=$1
  docker ps -a --format '{{.Names}}' | grep -Eq "^${container_name}\$"
}

# Function to exit the script if a container exists
check_and_exit_if_exists() {
  local container_name=$1
  if container_exists "$container_name"; then
    echo $'\xE2\x9C\x98' " Container $container_name already exists. Exiting script."
    exit 1
  fi
}

# Function to wait for a Redis service to be ready
wait_for_redis() {
  local container_name=$1
  local password=$2
  echo "Waiting for Redis to be ready..."
  echo "TRYING TO EXECUTE : docker exec ${container_name} redis-cli -a ${password} ping | grep -qP '^PONG'"
  until docker exec "$container_name" redis-cli -a "$password" ping | grep -qP '^PONG'; do
    sleep 2
  done
  echo $'\xE2\x9C\x93' " SUCCESSFULLY CREATING REDIS CONTAINER "
}

# Function to wait for a PostgreSQL service to be ready
wait_for_postgres() {
  local container_name=$1
  local db_user=$2
  local db_name=$3
  echo "Waiting for PostgreSQL to be ready..."
  echo "TRYING TO EXECUTE : docker exec ${container_name} pg_isready -U ${db_user} -d ${db_name}" 
  until docker exec "$container_name" pg_isready -U "$db_user" -d "$db_name"; do
    sleep 2
  done
  echo $'\xE2\x9C\x93'" SUCCESSFULLY CREATING POSTGRES CONTAINER"
}

# # Function to create a PostgreSQL database
create_database() {
  local container_name=$1
  local db_user=$2
  local db_name=$3

  echo "Creating database $db_name..."
  docker exec "$container_name" psql -U "$db_user" -c "CREATE DATABASE $db_name;"
  echo $'\xE2\x9C\x93' " SUCCESSFULLY CREATING DATABASE : " "$db_name" " ON POSTGRES CONTAINER"
}

execute_sql_file() {
  local container_name=$1
  local db_user=$2
  local db_name=$3
  local sql_file=$4

  echo "Executing SQL file $sql_file..."

  docker cp ./"$sql_file" "$container_name":/"$sql_file"
  echo docker cp ./"$sql_file" "$container_name":/"$sql_file"

  docker exec -it "$container_name" psql -U "$db_user" -d "$db_name" -f "./""$sql_file"
  echo docker exec -it "$container_name" psql -U "$db_user" -d "$db_name" -f "./""$sql_file"
  
  echo \xE2\x9C\x93 " SUCCESSFULLY EXECUTE " "$sql_file"
  echo
}

# Extract values from $docker_compose_yml using yq version 4
# shellcheck disable=SC2046
export $(grep -v '^#' $env | xargs)
DB_USER=${POSTGRES_USER}
DB_NAME=${POSTGRES_DATABASE}
POSTGRES_CONTAINER_NAME=$(yq e '.services.postgres.container_name' $docker_compose_yml)
REDIS_CONTAINER_NAME=$(yq e '.services.redis.container_name' $docker_compose_yml)
MAIN_APP=$(yq e '.services.app.container_name' $docker_compose_yml)
SQL_FILE="starter.sql"

echo

# Check the main container apps, if exist then break this order
check_and_exit_if_exists $MAIN_APP

# Check docker compose yml
if test -z "$docker_compose_yml"; then
    echo $'\xE2\x9C\x98' "The first parameter is empty, we need docker compose yml file"
else
    if [ -f "$docker_compose_yml" ]; then
      echo $'\xE2\x9C\x93' $docker_compose_yml " DETECTED"
    else
      echo $'\xE2\x9C\x98' $docker_compose_yml" NOT DETECTED"
      exit 1
    fi
fi

# Check docker compose yml
if test -z "$env"; then
    echo $'\xE2\x9C\x98' "The first parameter is empty, we need docker compose yml file"
else
    if [ -f "$env" ]; then
      echo $'\xE2\x9C\x93' $env " DETECTED"
    else
      echo $'\xE2\x9C\x98' $env" NOT DETECTED"
      exit 1
    fi
fi

# Check starter.sql
if [ -f "starter.sql" ]; then
  echo $'\xE2\x9C\x93' " starter.sql DETECTED"
else
  echo $'\xE2\x9C\x98' "starter.sql NOT DETECTED"
  exit 1
fi

echo

# Check and run Redis service
echo "READY TO COOK :" $REDIS_CONTAINER_NAME
if container_exists "$REDIS_CONTAINER_NAME"; then
  echo "-- Container $REDIS_CONTAINER_NAME already exists. Skipping Redis service."
else
  docker-compose -f $docker_compose_yml up --build -d redis
  wait_for_redis "$REDIS_CONTAINER_NAME" "$REDIS_PASSWORD"
fi
echo $'\xE2\x9C\x93' " COOKING" $REDIS_CONTAINER_NAME "SUCCESSFULLY !!!!"
echo

# Check and run PostgreSQL service
echo "READY TO COOK :" $POSTGRES_CONTAINER_NAME
if container_exists "$POSTGRES_CONTAINER_NAME"; then
  echo "-- Container $POSTGRES_CONTAINER_NAME already exists. Skipping PostgreSQL service."
else
  docker-compose -f $docker_compose_yml up --build -d postgres
  wait_for_postgres "$POSTGRES_CONTAINER_NAME" "$DB_USER" "$DB_NAME"
  create_database "$POSTGRES_CONTAINER_NAME" "$DB_USER" "$DB_NAME"
  execute_sql_file "$POSTGRES_CONTAINER_NAME" "$DB_USER" "$DB_NAME" "$SQL_FILE"
fi
echo $'\xE2\x9C\x93' " COOKING" $POSTGRES_CONTAINER_NAME "SUCCESSFULLY !!!!"
echo

# Start the golang_api service
echo "READY TO COOK THE MAIN KING"
docker-compose -f $docker_compose_yml up --build -d app
echo $'\xE2\x9C\x93' " COOKING THE MAIN KING SUCCESSFULLY !!!!"