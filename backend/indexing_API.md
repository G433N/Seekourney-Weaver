# REST API for indexers

**Note that `indexing` package implements all of the indexing side of the API.
A small amount of boiler-plate needs to be imported for it to run.
The purpose of this file is for documentation and for those who wish to
write an indexer in a language besides Go.**



The general standard for the API's requests and responses follow the
[jsend](https://github.com/omniti-labs/jsend) specifications which are very minimal.

HTTP response codes will be ignored, only JSON data will be checked.

`.json` ending in request path as stated by `jsend` will be omitted,
as responses will always be in json-format.

Note that JSON specifications do not support trailing commas.
Adding trailing commas in JSON responses may cause the parser to fail.

Only absolute file paths are used for requests.

While the main server starts up an indexer, it will check `stderr`.
If the initial ping-request is not able to be responded to,
an error message should be written to `stderr`.


## Ports

Allowed port numbers for indexers are any non-occupied port in the range [39 000, 39 499].

The port for a given indexer is calculated dynamically and sent as
command-line argument when starting up indexer.


## Possible requests and responses

A ping request will be send to the indexer shortly after attempted startup
as part of checking status.
```
GET /ping
```
Indexer must respond with:
```json
{
    "status" : "success",
    "data": {
        "message": "pong",
    }
}
```

Request for indexing of a folder or file takes the form:
```
POST /index/FILE_OR_FOLDERPATH
```
Indexer must respond with:
```json
{
    "status": "success",
    "data": {
        "message": "OPTIONAL MESSAGE",
    }
}
```
Alternatively, if indexing request cannot be handled:
```json
{
    "status": "fail",
    "data": {
        "message": "HUMAN-READABLE ERROR MESSAGE",
    }
}
```
The indexer then runs indexing on the given path and sends its own HTTP request
to Core:
```
POST /push/docs
to localhost port 8080
```
With request:
```json
{
    "status": "success",
    "data": {
        "documents": [
            {
                "path": "FILEPATH",
                "source": 0,
                "words": {
                    "SOMEWORD": 42,
                    "ANOTHER-WORD": 5,
                }
            },
            {
                "path": "PATHFORSOMEWEB",
                "source": 1,
                "words": {
                    ...
                }
            },
            ...
        ],
    }
}
```
Alternatively, if indexing failed:
```json
{
    "status": "fail",
    "data": {
        "message": "HUMAN-READABLE ERROR MESSAGE",
    }
}
```

**`"source"` value in response must be `0` for a local file,
or `1` for web file (e.g. HTML).**

Note that if none or only a single document was produced the `"documents"`
field should still be present.
For consistency, an array with the same key is
still used in the response when indexing a single file.


A shutdown request may be sent to the indexer from Core (main server).
```
GET /shutdown
```
Indexer must respond with:
```json
{
    "status": "success",
    "data": {
        "message": "exiting",
    }
}
```
And immediately exit all it's associated processes.
