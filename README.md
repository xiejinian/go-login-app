# Go Login Application

A secure web application built with Go, featuring user authentication and session management.

## Features

- User Authentication
- Session Management
- Secure Password Storage
- Dashboard Access
- GORM Database Integration
- Gin Web Framework

## Technologies Used

- Go (Golang)
- Gin Web Framework
- GORM (ORM)
- SQLite Database
- bcrypt for Password Hashing

## Prerequisites

- Go 1.16 or higher
- SQLite

## Installation

1. Clone the repository:
```bash
git clone [your-repository-url]
cd go-login-app
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The application will be available at `http://localhost:8080`

## Default Login Credentials

- Username: admin
- Password: password123

## Project Structure

- `main.go`: Server initialization and route setup
- `handlers/`: Authentication and request handlers
- `database/`: Database operations and initialization
- `models/`: Data models
- `templates/`: HTML templates
- `static/`: Static files (CSS, JS, etc.)

## Security Features

- Password hashing with bcrypt
- Secure session management
- HttpOnly cookies
- Input validation
- Error handling

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
