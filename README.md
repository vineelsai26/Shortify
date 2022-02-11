# Shortify

Url Shortener webapp with java script

## Download

Download docker container using command

```sh
docker pull ghcr.io/vineelsai26/shortify:main
```

### Run

Run the container on port 5000 and replace MONGODB_URL with url to mongodb database

```sh
docker run -p 5000:5000 -e MONGODB='MONGODB_URL' -d ghcr.io/vineelsai26/shortify
```
