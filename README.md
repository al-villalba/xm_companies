
# Technical Test for XM by Alvaro Villalba

## Code structure

The micreservice can be found in the folder `producer`. As required, it 
implements a REST API with the operations *create*, *patch*, *delete*, and
*get* (one). Definition of routes for this operations are found in
`producer/router.go`.

Everything starts in `producer/main.go`. When the application starts, dependent
services are loaded, reporting error or panicking if loading fails.

I use mysql as Db and Kafka for events. Both services as well as the application
are dockerised. The configuration files are located in `compose.yml` and
`producer/Dockerfile`.

The action, as described in the diagram `xm-tech-test.drawio.svg`, goes first
to the database. If there are no error, then an envent is produced in kafka
(only for muting operations).

## Setup

Everything has been prepared for the application and services to start by
calling the command

```sh
docker compose up
```

Despite the application checks the dependent services and retry connecting to
them several times, it may give up if these take too long to start. If this
case, it's worth to retry `docker compose up`. Services start quicker when
containers are cached. If still the application gives up before the services
start, you can try an alternative setup. See further below

## How to use it

Docker compose up will initialize the database, create the table `companies`,
and create the topic `xm-companies` in kafka. In case they are not, they can be
created manually (see below). Following are the curl commands to hit the API:

* Create:
```sh
curl -i -X POST -H 'Content-Type: application/json' \
-d '{"id": "41931b35-0541-4201-99b7-4d61026e0941", "name": "test company 2", "description": "", "amt_employees": 4000, "registered": true, "type": "Sole Proprietorship"}' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwcHVzZXIiLCJyb2xlIjoid3JpdGVyIn0.Oz6R6Kq8DkDTg6FqJY-Fw3kI0_XJh8ATKP9x3p-7OQ4' \
http://localhost:3000/company
```

* Patch:
```sh
curl -i -X PUT -H 'Content-Type: application/json' \
-d '{"id": "41931b35-0541-4201-99b7-4d61026e0941", "name": "test company 2", "description": "", "amt_employees": 4001, "registered": true, "type": "Sole Proprietorship"}' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwcHVzZXIiLCJyb2xlIjoid3JpdGVyIn0.Oz6R6Kq8DkDTg6FqJY-Fw3kI0_XJh8ATKP9x3p-7OQ4' \
http://localhost:3000/company
```
* Get:
```sh
curl -i http://localhost:3000/company/41931b35-0541-4201-99b7-4d61026e0941
```

* Delete:
```sh
curl -i -X DELETE -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwcHVzZXIiLCJyb2xlIjoid3JpdGVyIn0.Oz6R6Kq8DkDTg6FqJY-Fw3kI0_XJh8ATKP9x3p-7OQ4' \
http://localhost:3000/company/41931b35-0541-4201-99b7-4d61026e0941
```

A JWT token has been generated for this test. This has be done with
`GenerateJwtTokenString`, found in `common/common/jwtAuth.go`

The events can be consumed in the kafka container with 
`/opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic xm-companies`

## Alternative Setup

1. Start the services: `docker compose up mysql zookeeper kafka`

2. Once they have started, in the `producer` folder edit .env the IP addresses
of containers kafka and mysql. Then you can start the application by running
`source .env.sh && go run .`

## Manual Setup for MySQL and Kafka

1. The sql script to create the table is located in `producer/data/migration-init.sql`
Enter the mysql container: `docker exec -it mysql /bin/sh` and run the script
`mysql -u xm -p xm_tech_test < /docker-entrypoint-initdb.d/migration-init.sql`
(the password is in the configuration files)

2. To create the topic in kafka, enter the container `docker exec -it karka /bin/sh`
and run `/opt/kafka/bin/kafka-topics.sh --zookeeper zookeeper:2181 --create --topic xm-companies --partitions 1 --replication-factor 1`
