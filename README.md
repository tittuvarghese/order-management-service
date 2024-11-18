# Order Management Service for E-Commerce Application

[![GoDoc](https://pkg.go.dev/badge/github.com/tittuvarghese/ss-go-order-management-service)](https://pkg.go.dev/github.com/tittuvarghese/ss-go-order-management-service)
[![Build Status](https://travis-ci.org/tittuvarghese/ss-go-order-management-service.svg?branch=main)](https://travis-ci.org/tittuvarghese/ss-go-order-management-service)

The **Order Management Service** is a microservice responsible for handling customer orders in an e-commerce platform. It provides gRPC endpoints to manage orders, including creating new orders, retrieving orders by customer or order ID, updating order status, and more.

This service is integrated with other microservices like **Product Service** and **Customer Service** for complete order processing functionality.

## API Overview

The **Order Management Service** exposes the following gRPC methods for managing orders:

### 1. **Create Order**
- **RPC Method**: `CreateOrder`
- **Request Type**: `CreateOrderRequest`
- **Response Type**: `CreateOrderResponse`
- **Description**: Creates a new order for the specified customer with order items and shipping details.

#### Request (CreateOrderRequest)
```proto
message CreateOrderRequest {
  string customer_id = 1;
  repeated OrderItem items = 2;
  Address Address = 3;
  string Phone = 4;
}
```

#### Response (CreateOrderResponse)
```proto
message CreateOrderResponse {
  string Message = 1; // Success or failure message
  Order Order = 2;    // The created order
}
```

### 2. **Get Orders**
- **RPC Method**: `GetOrders`
- **Request Type**: `GetOrdersRequest`
- **Response Type**: `GetOrdersResponse`
- **Description**: Retrieves all orders associated with a specific customer.

#### Request (GetOrdersRequest)
```proto
message GetOrdersRequest {
  string customer_id = 1;
}
```

#### Response (GetOrdersResponse)
```proto
message GetOrdersResponse {
  string Message = 1;  // Success or failure message
  repeated Order Orders = 2;  // List of orders for the specified customer
}
```

### 3. **Get Order by ID**
- **RPC Method**: `GetOrder`
- **Request Type**: `GetOrderRequest`
- **Response Type**: `GetOrderResponse`
- **Description**: Retrieves a single order by its order ID.

#### Request (GetOrderRequest)
```proto
message GetOrderRequest {
  string customer_id = 1; // Customer ID
  string order_id = 2;    // Order ID
}
```

#### Response (GetOrderResponse)
```proto
message GetOrderResponse {
  string Message = 1;   // Success or failure message
  Order order = 2;      // The requested order details
}
```

### 4. **Update Order Status**
- **RPC Method**: `UpdateOrderStatus`
- **Request Type**: `UpdateOrderStatusRequest`
- **Response Type**: `UpdateOrderStatusResponse`
- **Description**: Updates the status of an existing order (e.g., from "PENDING" to "SHIPPED").

#### Request (UpdateOrderStatusRequest)
```proto
message UpdateOrderStatusRequest {
  string order_id = 1;  // The order ID
  string customer_id = 2;  // Customer ID
  string status = 3;    // New status for the order
}
```

#### Response (UpdateOrderStatusResponse)
```proto
message UpdateOrderStatusResponse {
  string Message = 1; // Success or failure message
}
```

### 5. **Cancel Order (Currently Commented Out)**
- **RPC Method**: `CancelOrder`
- **Request Type**: `CancelOrderRequest`
- **Response Type**: `CancelOrderResponse`
- **Description**: Cancels an order (currently not implemented).

#### Request (CancelOrderRequest)
```proto
message CancelOrderRequest {
  string customer_id = 1;
  string order_id = 2;
}
```

#### Response (CancelOrderResponse)
```proto
message CancelOrderResponse {
  string Message = 1; // Success or failure message
}
```

## Order and OrderItem Message Definitions

The **Order** and **OrderItem** message structures are defined as follows:

```proto
message OrderItem {
  string product_id = 1; // Product ID
  int32 quantity = 2;    // Quantity of the product
  double price = 3;      // Price of the product
}

message Order {
  string order_id = 1;      // Unique order ID
  string customer_id = 2;   // Customer ID
  repeated OrderItem items = 3;  // List of items in the order
  double total_price = 4;   // Total price for the order
  string status = 5;        // Current order status
  Address Address = 6;      // Shipping address
  string Phone = 7;         // Customer phone number
}
```

### Address Message
```proto
message Address {
  string AddressLine1 = 1;  // Street address line 1
  string AddressLine2 = 2;  // Street address line 2 (optional)
  string City = 3;          // City
  string State = 4;         // State/Province
  string Country = 5;       // Country
  string Zip = 6;           // Zip or postal code
}
```

### Order Status Enum
The **OrderStatus** enum defines the possible states an order can be in:

```proto
enum OrderStatus {
  PENDING = 0;
  PROCESSING = 1;
  SHIPPED = 2;
  DELIVERED = 3;
  CANCELED = 4;
  RETURNED = 5;
}
```

## Running the Service Locally

### Prerequisites

Before running the Order Management Service locally, ensure the following:

- Go 1.18 or higher
- Protocol Buffers (Protobuf) Compiler (`protoc`)
- gRPC Go Plugin for Protobuf (`protoc-gen-go` and `protoc-gen-go-grpc`)

### Steps to Run Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/tittuvarghese/ss-go-order-management-service.git
   cd order-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate gRPC code from the `proto` file:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/order.proto
   ```

4. Start the Order Service:
   ```bash
   go run cmd/main.go
   ```

The service will start and listen for gRPC requests on the specified port (e.g., `50051`).


## Example Usage with Gateway Service

### Create a New Order
```bash
curl -X POST http://localhost:8080/order/create \
   -d '{"customer_id": "12345", "items": [{"product_id": "67890", "quantity": 2, "price": 49.99}], "Address": {"AddressLine1": "123 Main St", "City": "Anytown", "State": "CA", "Country": "US", "Zip": "12345"}, "Phone": "555-1234"}' \
   -H "Content-Type: application/json"
```

### Get All Orders for a Customer
```bash
curl -X GET http://localhost:8080/orders \
   -d '{"customer_id": "12345"}' \
   -H "Authorization: Bearer <your_jwt_token>"
```

### Get Order Details by ID
```bash
curl -X GET http://localhost:8080/order/12345 \
   -d '{"customer_id": "12345", "order_id": "67890"}' \
   -H "Authorization: Bearer <your_jwt_token>"
```

### Update Order Status
```bash
curl -X POST http://localhost:8080/order/update-status \
   -d '{"order_id": "67890", "customer_id": "12345", "status": "SHIPPED"}' \
   -H "Authorization: Bearer <your_jwt_token>"
```

## Architecture

The **Order Management Service** is an integral part of the e-commerce ecosystem. It communicates with other services (like the **Customer Service** and **Product Service**) through gRPC to manage order-related tasks efficiently.

- **gRPC Communication**: The service uses gRPC for high-performance, low-latency communication between microservices.
- **Order Lifecycle**: Handles the creation, updating, and retrieval of orders, including managing the status of each order (e.g., "PENDING", "SHIPPED").
- **Database**: Stores order details, customer information, and order status.
- **JWT Authentication**: Secures endpoints by using JWT for authenticated access.
