# uTurn - A Simple URL Shortener

Deploy In Two steps.

1. Clone the Repository
2. Run it
```bash
$ go run main.go
```

## Endpoints

- GET `/ping`

  Health Check, gives you a 200 if all is okay

- GET `/urls`

  Gives you a list of all the URLs stored

- POST `/urls`

  Accepts - 
  ```json
  {
    "Shortcode": "ABCDE", # this is optional, we'll generate one if you don't have one
    "URL": "www.my_funny_website.com"
  }
  ```

- GET `/<shortcode>`
  
  For example, GET `http://localhost:8000/abc` should take you to Google!
  Returns an HTTP 304 to the URL Stored at that Shortcode.
