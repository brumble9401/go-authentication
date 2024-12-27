# Golang Authentication Project

## Introduction

This project is a simple authentication system implemented in Golang. It demonstrates how to handle user registration, login, and session management using JWT (JSON Web Tokens). It supports both username/password login and Google login.

## Tech Stack

- **Golang**: The main programming language used for developing the authentication system.
- **ScyllaDB**: A high-performance NoSQL database used for storing user data and authentication information.
- **Redis**: An in-memory data structure store used for session management.

## Features

- User Registration
- User Login
- Google Login
- JWT-based Authentication
- Secure Password Storage
- Session Management

## Prerequisites

- Go 1.16 or higher
- A running instance of ScyllaDB
- A running instance of Redis

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
   # Update .env with your ScyllaDB, Redis credentials, and Google OAuth credentials
   ```

4. Run the application:

   ```sh
   make run
   ```

5. Run database migrations:

   ```sh
   make migrate-up
   ```

## Usage

### Register a new user

Send a POST request to `/register` with the following JSON payload:

```json
{
  "username": "yourusername",
  "password": "yourpassword",
  "email": "youremail@example.com",
  "full_name": "Your Full Name",
  "role": "USER"
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

### Google Login

Send a GET request to `/auth/google/login` to initiate the Google login process. After successful authentication, you will receive a JWT token in response.

### Access protected routes

Include the JWT token in the `Authorization` header as follows:

```sh
Authorization: Bearer <your-jwt-token>
```

## Project Structure

```
/golang-authentication
|-- /api
|-- /bin
|-- /cmd
|-- /config
|-- /interfaces
|-- /middleware
|-- /migrations
|-- /models
|-- /querybuilder
|-- /redis
|-- /repository
|-- /scylla
|-- /services
|-- /utils
|-- .env.example
|-- .gitignore
|-- go.mod
|-- go.sum
|-- Makefile
|-- README.md
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

## Contact

For any inquiries, please contact [your-email@example.com].
