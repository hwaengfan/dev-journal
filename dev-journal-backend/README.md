To run this locally, make sure to set up a MySQL database server

#### Create a .env file containing:

```
PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_ADDRESS=127.0.0.1
DB_NAME=dev-journal-database
```

To run the server: `make run`

#### Create a Docker container for the MySQL database:

```
docker run --name dev-journal-database -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql
docker start dev-journal-database
```

#### Migrate the tables for the MySQL database:

```
make migrate-up
```

#### Clear tables in the MySQL database:

```
make migrate-down
```
