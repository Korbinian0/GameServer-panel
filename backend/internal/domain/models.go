package domain

type User struct {
    ID       string
    Email    string
    Password string
    Role     string
    CreatedAt string
}

type Role struct {
    Name        string
    Permissions []string
}

type Node struct {
    ID          string
    Platform    string
    Hostname    string
    IPAddress   string
    Capabilities []string
    Version     string
    LastSeen    string
}

type ServerInstance struct {
    ID          string
    NodeID      string
    GameType    string
    Status      string
    PID         string
    UptimeSeconds int64
    PlayerCount int
    MaxPlayers  int
}
