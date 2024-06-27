To run this locally, make sure to create a Docker container for a MySQL server and set up the credentials in the environment variables.

#### Create a Docker container:

```
docker run --name dev-journal-database -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql
docker start dev-journal-database
```

#### Create a .env file containing:

```
PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_ADDRESS=127.0.0.1
DB_NAME=dev-journal-database
```

To run the server: `make run`
