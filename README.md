##### JSON API is base bank system
##### Features
- login
- jwt auth
- transfer

##### Usage
Run docker postgres image
```bash
docker run --name some_postgres -e POSTGRES_PASSWORD=go_bank -p 5432:5432 -d postgres
```
Run application with seed 
```bash
make seed
```
This seed creates user with the number 7625482 and password "secret" with a balance 50000

Send POST request to localhost:8080/login  with  seedable account
```json
{
  "number": 7625482,
  "password": "secret"
}
```
The response contains a 'token' field that should be inserted into the header for subsequent requests:
header: x-jwt-token
value: {token}

Create your account:
POST localhost:8080/account
```json
{
  "firstName": "MyName",
  "lastName": "MySecondName",
  "password": "mysecretpass"
}
```
List all accounts:
GET localhost:8080/account

Example transfer:
POST localhost:8080/transfer
Note: you need to create a second account to transfer to
```json
{
  "fromAccount": 7625482,
  "toAccount" : 7638280, // type here your created account number
  "amount": 2000
}
```

