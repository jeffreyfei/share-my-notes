# Share  My Notes
- Realtime note sharing platform

# Development
### Prerequisite
1. Clone the repo to `$GOPATH/src`

1. Install Go

1. Install Glide

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
### Get dependencies
```
bash dev.sh install
```
### Build server
```
bash dev.sh build
```
### Run server
```
bash dev.sh run-server
```