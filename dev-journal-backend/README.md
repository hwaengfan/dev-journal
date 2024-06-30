## To run this locally, make sure to:

1. Set up a MySQL database server
2. Fill in the correct environment variables
3. Migrate the tables

#### Create a .env file containing:

```
PUBLIC_HOST=http://localhost
PORT=8080

DB_USER=root
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=dev-journal-database

JWT_EXPIRATION_IN_SECONDS=86400
JWT_SECRET=dev-journal-secret
```

#### Create and run a Docker container for the MySQL database server:

```
docker run --name dev-journal-database -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql
docker start dev-journal-database
```

#### Create a dev-journal-database schema in the server by connecting it via CLI or through an editor and run:

```
CREATE SCHEMA 'dev-journal-database';
```

#### Migrate the tables to the MySQL database:

```
make migrate-up
```

#### Run the backend server:

```
make run
```

#### Clear tables in the MySQL database:

```
make migrate-down
```
