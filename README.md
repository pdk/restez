# restez

handling REST more EZ

    curl -s -S -X POST http://localhost:8080/new -d '{"name": "Bob", "age": 45, "favoriteFood": "ice cream"}' | jq
    curl -s -S -X POST http://localhost:8080/new -d '{"name": "Marty", "age": 13, "favoriteFood": "grapes"}' | jq

    curl -s -S -X GET http://localhost:8080/list | jq
