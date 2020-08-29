package config


type Config struct{
    User string `json:"user"`
    SMTPServerAddr string `json:"server_addr"`
    Passwd string `json:"passwd"`
}
