# Microservices Project

This project implements a set of microservices for managing accounts, orders, and catalogs. The services are built using Go and gRPC, with a GraphQL API Gateway to provide a unified interface for interacting with the microservices. The services use different databases for persistence:

- **Account and Order Services**: PostgreSQL
- **Catalog Service**: Elasticsearch

## Project Structure

The project consists of the following main components:

- **Account Service**: Manages user accounts.
- **Catalog Service**: Manages products and catalogs.
- **Order Service**: Manages customer orders.
- **GraphQL API Gateway**: Provides a unified GraphQL interface for interacting with the above services.

## Getting Started

Follow the steps below to get the project running locally.

### 1. Clone the repository


git clone <repository-url>
cd <project-directory>

## 2. Start the services using Docker Compose

Ensure you have Docker and Docker Compose installed. Then, run the following command to build and start the services in the background:

docker-compose up -d --build


## 3. Access the GraphQL Playground

Once the services are running, you can access the GraphQL playground at:

[http://localhost:8000/playground](http://localhost:8000/playground)

## gRPC File Generation

To generate the necessary gRPC files, follow these steps:

1. Download the Protocol Buffers compiler (`protoc`):

   Download the file from the specified link, unzip it, and move the `protoc` binary to `/usr/local/bin/`.

2. Install the necessary Go tools for gRPC:

   Install `protoc-gen-go` and `protoc-gen-go-grpc` tools using Go.

3. Update your `PATH` to include the Go binaries.

4. Create the `pb` folder in your project directory.

5. In the `account.proto` file, add the following line to set the Go package:

option go_package = "./pb";


6. Finally, run the `protoc` command to generate the gRPC files.

## GraphQL API Usage

The GraphQL API provides a unified interface to interact with all the microservices.

### Query Accounts

query { accounts { id name } }

### Create an Account

mutation { createAccount(account: {name: "New Account"}) { id name } }

### Query Products

query { products { id name price } }

### Create a Product

mutation { createProduct(product: {name: "New Product", description: "A new product", price: 19.99}) { id name price } }<D-z>u

### Create an Order

mutation { createOrder(order: {accountId: "account_id", products: [{id: "product_id", quantity: 2}]}) { id totalPrice products { name quantity } } }

### Query Account with Orders

query { accounts(id: "account_id") { name orders { id createdAt totalPrice products { name quantity price } } } }

## Advanced Queries

### Pagination and Filtering

query { products(pagination: {skip: 0, take: 5}, query: "search_term") { id name description price } }

### Calculate Total Spent by an Account

query { accounts(id: "account_id") { name orders { totalPrice } } }
