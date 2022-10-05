## Tech Stack


- [Express (Server Framework) ](https://expressjs.com/)
- [morgan (Express Logging Framework) ](https://www.npmjs.com/package/morgan)

- [tsoa (Build OpenAPI-compliant REST APIs using TypeScript and Node) ](https://github.com/lukeautry/tsoa)
- [Swagger UI ( Rest  API explorer )](https://www.npmjs.com/package/swagger-ui)


## Server URL's
- [Swagger UI ] (http://localhost:3000/docs/)
- [Get User By Id 1 ]  (http://localhost:3000/users/1)

## Build from source

1. Install dependencies.

   ```sh
   yarn install
   ```

2. Build the production server.

   ```sh
   yarn build
   ```

3. Run the server.
   ```sh
   yarn start
   ```

## Run tests

```sh
yarn test
```

### Run Client Commands

```sh
curl -X GET http://localhost:3000/users
```


```sh
curl -X GET http://localhost:3000/users/1
```

```sh
curl -X 'POST' \
  'http://localhost:3000/users' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "Shashank",
  "role": "Admin"
}'
```

```sh
curl -X 'DELETE' \
  'http://localhost:3000/users/2' \
  -H 'accept: application/json'
```