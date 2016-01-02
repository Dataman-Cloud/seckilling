# omega-god
omega admin platform

## Setup

Clone the omega-god repo
```
# git clone git@github.com:Dataman-Cloud/omega-god.git
```

Start the containers
```
# make init
```

Create the database omega-god and omega
```
# make create-db
```

Migrate the database omega-god and import the testing data in database omega
```
# make init-db
```

Rerstart the omega-god after the database is ready
```
# make restart-god
```

## Clean up the ENV
```
# make cleanup
```
