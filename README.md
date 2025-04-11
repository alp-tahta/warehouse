# warehouse
Order-Based Shelf Management System

## How to run project
[Click here to how to run and learn how does it works](https://youtu.be/x0Ft0-B7ak8)
```sh
make up
```
Or if you dont have make
```sh
docker compose up -d
```

### Check api health
```sh
curl --location 'localhost:8080/health'
```

### Create order for user
```sh
curl --location 'localhost:8080/order' \
--header 'Content-Type: application/json' \
--data '{
  "customer_id": "13214459-18dd-48f0-ae51-ed42d640a172",
  "order_items": [
    {
      "product_id": 1001,
      "quantity": 2
    },
    {
      "product_id": 1002,
      "quantity": 1
    }
  ]
}'
```

### ReadBarcode example at warehouse
```sh
curl --location --request PUT 'localhost:8080/barcode/13214459-18dd-48f0-ae51-ed42d640a172*ae306152-7872-4f48-a82f-86c8b2824e8e*1002'
```

### Get Shelves' Details as RESTful-API
```sh
curl --location 'localhost:8080/shelf'
```

### Visit [Shelves' Details](http://localhost:8080/shelf-html) for UI

## Open Shell into docker container and connect db warehouse to see database tables
```bash
docker exec -it postgres psql -U postgres -d warehouse
```
then
```sql
SELECT * FROM orders;
SELECT * FROM order_items;
SELECT * FROM barcodes;
SELECT * FROM shelves;
```

## Improvements TODO
- This is monolit service, can be turn into microservice architecture.
- Some DB jobs can turn in transactional for reliablity.


mockgen -source=/home/alp-tahta/projects/warehouse/internal/barcode/barcodeinterface.go -destination=/home/alp-tahta/projects/warehouse/internal/barcode/mock_barcoder.go -package=barcode
mockgen -source=/home/alp-tahta/projects/warehouse/internal/repository/repository.go -destination=/home/alp-tahta/projects/warehouse/internal/repository/mock_repository.go -package=repository
mockgen -source=/home/alp-tahta/projects/warehouse/internal/service/service.go -destination=/home/alp-tahta/projects/warehouse/internal/service/mock_service.go -package=service