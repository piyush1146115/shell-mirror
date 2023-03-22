# shell-mirror

```bash
curl --request POST \
  --url http://localhost:8088/execute \
  --header 'content-type: application/json' \
  --data '{"command": "pwd"}'
```

```bash
curl --request POST \
  --url http://localhost:8088/execute \
  --header 'content-type: application/json' \
  --data '{"command": "ls -l"}'
```

```bash
curl --request POST \
  --url http://localhost:8088/execute \
  --header 'content-type: application/json' \
  --data '{"command": "cat xyz.txt"}'
```