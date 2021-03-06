# Share  My Notes
- Realtime note sharing platform

# Development
- Default development server route `localhost:3000`
### Prerequisite
1. Clone the repo to `$GOPATH/src`

1. Install Go

1. Install Glide

1. Setup Google credentials in your Google Account
```
export GOOGLEKEY=<Your google client ID>
export GOOGLESECRET=<Your google client secret>
```

### Setup Postgres
1. Follow instructions on https://www.postgresql.org/ to install Postgresql 10 on your system
1. Initialize Postgres
```
/usr/pgsql-10/bin/postgresql-10-setup initdb
systemctl enable postgresql-10
systemctl start postgresql-10

service postgresql initdb
chkconfig postgresql on
```
3. Setup local connecction to Postgres in http://suite.boundlessgeo.com/docs/latest/dataadmin/pgGettingStarted/firstconnect.html
- In pg_hba.conf change peer to trust (Applicable in dev enviroment only!)

4. Create database `"development"` and `"testing"`
### Get dependencies
```
bash dev.sh up
```
### Build server
```
bash dev.sh build
```
### Run load balancer
```
bash dev.sh run-load-balancer
```
### Run server
```
bash dev.sh run-server
```
### Run tests
```
bash dev.sh test
```