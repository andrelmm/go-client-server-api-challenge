# Go Client Server Api Challenge

This project demonstrates a simple client-server interaction in Go to fetch exchange rates and store them in a SQLite database.

The project consists of two main components:
- `client.go`: A client application that fetches the current USD to BRL exchange rate from the server and saves it to a file.
- `server.go`: A server application that fetches the current exchange rate from an external API, stores it in a SQLite database, and provides it to clients upon request.

### Prerequisites

- Go (Golang)
- SQLite


### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/exchange-rate-server.git
   cd go-client-server-api-challenge
   ```

2. Install dependecies:

   ```bash
   go mod tidy
   ```

3. Run the server:

   ```bash
   go run server.go
   ```

5. In a separate terminal, run the client:

   ```bash
   go run client.go
   ```

### Project Structure

The project follows this structure:

- `client.go`: The client application code that fetches the exchange rate from the server and saves it to a file.
- `server.go`: The server application code that fetches the exchange rate from an external API, stores it in the database, and serves client requests.
- `cotation.go`: The SQLite database file where the exchange rates are stored.
- `cotacao.txt`: The text file where the fetched exchange rate is saved by the client.
- `README.MD`: This documentation file.


