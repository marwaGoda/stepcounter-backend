# stepcounter-backend

The Company Steps Leaderboard Backend is a RESTful API built to manage teams and their step counters. It allows you to create teams, add users to teams, increment step counters for users, and retrieve team and user information.

## Getting Started

These instructions will help you set up and run the backend on your local machine.

### Prerequisites

You will need the following tools installed on your system:

- Go (https://golang.org/)
- Git (https://git-scm.com/)

### Installation

1. Clone the repository:

```
git clone https://github.com/marwaGoda/stepcounter-backend.git
```

2. Change into the project directory:

```
cd backend
```

3. Install the required dependencies:

```
go mod download
```

4. Run the backend server:

```
go run main.go
```

The server will start running on `http://localhost:8080`.

## API Endpoints

The API provides the following endpoints:

- `POST /teams`: Create a new team.
- `GET /teams`: Retrieve a list of all teams.
- `GET /teams/{teamID}`: Retrieve the details of a specific team.
- `POST /teams/{teamID}/users`: Add a new user to a team.
- `GET /teams/{teamID}/users`: Retrieve all users in a team.
- `POST /teams/{teamID}/users/{userID}/counters`: Increment the step count for a user in a team.
- `GET /teams/{teamID}/counters`: Retrieve the current total steps taken by a team.
- `POST /teams/{teamID}/users/{userID}/counters/increment`: Increment the step counter of a user in a team by a value.



## Running Tests

To run the unit tests for the backend, use the following command:

```
go test ./api
```

## Built With

- Go - The programming language used.
- Gorilla Mux - The HTTP router used for handling requests.
