```go
package main

type Client struct {
    ID    string
    Name  string
    Email string
    Phone string
}

type Project struct {
    ID          string
    ClientID    string
    Title       string
    Description string
    Status      string // "active", "completed"
}

type Communication struct {
    ID        string
    ClientID  string
    Subject   string
    Content   string
    Timestamp string
}

type Invoice struct {
    ID        string
    ProjectID string
    Status    string // "draft", "sent", "paid"
    Amount    float64
}

type Activity struct {
    ID        string
    ClientID  string
    Type      string // "communication", "milestone", "note"
    Title     string
    Content   string
    Timestamp string
}
```