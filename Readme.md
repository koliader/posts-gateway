# Posts gateway

## Installation

Install docker and wsl (if you on windows)
Copy two services

- [Auth service](https://github.com/koliader/posts-auth)
- [Posts service](https://github.com/koliader/posts-post)

### !!! THIS SERVICES SHOULD BE IN THE SAME DIRECTORY AS THIS GATEWAY !!!

```
$ docker compose up
```

## Usage

Service runs on port 8080 (specified in example.env)
Before usage set the TOKEN_KEY in example.env
The TOKEN_KEY should be the same as in the auth service
