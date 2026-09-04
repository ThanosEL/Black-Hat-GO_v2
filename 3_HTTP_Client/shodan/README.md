Here is a high-level overview of the typical steps for preparing and building an API client:
1. Review the service’s API documentation
2. Design a logical structure for the code in order to reduce complexity and repetition
3. Define request or response types, as necessary, in Go
4. Create helper functions and types to facilitate simple initialization, authentication, and communication to reduce verbose or repetitivelogic.
5. Build the client that interacts with the API consumer functions and types

### Designing the Project Structure
```
tree .
.
├── cmd
│   └── shodan
│       └── main.go
├── README.md
└── shodan
    ├── api.go
    ├── host.go
    └── shodan.go
```

- The main.go file defines package main and is used primarily as a con-sumer of the API you’ll build
- The files in the shodan directory—api.go, host.go, and shodan.go—define package shodan, which contains the types and functions necessary for com-munication to and from Shodan