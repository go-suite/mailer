# Mailer
A simple API mail server

## Usage
Mailer can be run in 2 different modes.

It can be open for every one sending a request or it can be limited to a list of users defined in `assets\users.toml` or `assets\users.yaml`

### Define a list of users allowed to use Mailer
`TOML`
```toml
[[users]]
username = "alice"
password = "alice_password"

[[users]]
username = "bob"
password = "bob_password"
```

`YAML`
```yaml
users:
  - username: "martin"
    password: "martin_password"
    authentication:
      server: "smtp.example.com"
      port: 587
      user:  "user"
      password: "123456"
  - username: "daniel"
    password: "daniel_password"
```

## Request
### Check if the server is online
`GET http://localhost:8081/check`

### Send an email
`POST http://localhost:8081/send`

```json
{
    "message": {
        "from": "your-email-address",
        "to": "first-email-address,second-email-address",
        "subject": "email-subject",
        "body": "html-email-body-content"
    },
    "authentication": {
        "server": "your-smtp-server",
        "port": server-port,
        "user": "smtp-user",
        "password": "smtp-user-password"
    }
}
```



### Send an email using an existing user

The authentication can be omitted if defined in the users list

`POST http://localhost:8081/send`

```json
{
	"message": {
		"from": "your-email-address",
		"to": "first-email-address,second-email-address",
		"subject": "email-subject",
		"body": "html-email-body-content"
	}
}

```
