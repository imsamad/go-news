# Go News

A web-based news application where users can create, delete, and edit news, while the admin has full control over managing all news posts. The system uses email-based authentication.

### Setup guide

After cloning the repo
#### Spin up MySql container

```sh
make docker
or
docker compose up -d
```
Note: MySQL container make take around 30sec for init, please wait for 30sec then run next cmds

#### Seed the db
```sh
make seed 
or
cd seed && go run .
```

#### Launch the app
```sh
cd app && go run .
```

Application would be up and running on port 3000

Test credentials

```js
Users: (user1, user2, user3, user4)@email.com,
Admin: admin@gmail.com 
Password: 123456
```
