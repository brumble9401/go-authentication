# Golang Authentication Project

## Introduction

This project is a simple authentication system implemented in Golang. It demonstrates how to handle user registration, login, and session management using JWT (JSON Web Tokens).

## Tech Stack

- **Golang**: The main programming language used for developing the authentication system.
- **ScyllaDB**: A high-performance NoSQL database used for storing user data and authentication information.

## Features

- User Registration
- User Login
- JWT-based Authentication
- Secure Password Storage

## Prerequisites

- Go 1.16 or higher
- A running instance of a database (e.g., PostgreSQL)

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/golang-authentication.git
   cd golang-authentication
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Set up environment variables:

   ```sh
   cp .env.example .env
   # Update .env with your database credentials and JWT secret
   ```

4. Run the application:
   ```sh
   go run main.go
   ```

## Usage

### Register a new user

Send a POST request to `/register` with the following JSON payload:

```json
{
  "username": "yourusername",
  "password": "yourpassword"
}
```

### Login

Send a POST request to `/login` with the following JSON payload:

```json
{
  "username": "yourusername",
  "password": "yourpassword"
}
```

You will receive a JWT token in response.

### Access protected routes

Include the JWT token in the `Authorization` header as follows:

```sh
Authorization: Bearer <your-jwt-token>
```

## Project Structure

```
/golang-authentication
|-- /controllers
|-- /models
|-- /routes
|-- /utils
|-- main.go
|-- .env.example
|-- go.mod
|-- go.sum
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

## Contact

For any inquiries, please contact [your-email@example.com].
