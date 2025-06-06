# REST API for indexers

**Note that `indexing` package implements all of the indexing side of the API.
A small amount of boiler-plate needs to be imported for it to run.
Only the step of converting a given file to a text string,
needs to be provided (if the indexer is indexing local files).
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


## Ports

Allowed port numbers for indexers are any non-occupied port in the range [39 000, 39 499].

The port for a given indexer is calculated dynamically and sent as
command-line argument when starting up the indexer.


## Possible requests and responses

When registering a new indexer, a request to get the name of the indexer
used to display to users on fronend clients will be sent.
```
GET /name
```
Indexer must respond with:
```
INDEXERNAME
```
as a plain-text string.

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
with body:
```json
POST /index
{
    "path":           "PATHTODATA", // E.g. filepath
    "type":           0,            // Source type, e.g. local file
    "collectionid":   "102983472",        // This is a generated hash
    "recursive":      true,         // Bool, the indexer defines what this means
    "parrallel":      true,         // Bool, the indexer defines what this means
}
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
to the main server:
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
                "collectionid": "102983472",
                "words": {
                    "SOMEWORD": 42,
                    "ANOTHER-WORD": 5,
                }
            },
            {
                "path": "PATHFORSOMEWEBSITE",
                "source": 1,
                "collectionid": "102983472",
                "words": {
                    ...
                }
            },
            ...
        ],
    }
}
```

**`"source"` value in response must be `0` for a local file,
or `1` for web file (e.g. HTML).**

Note that if none or only a single document was produced the `"documents"`
field should still be present.
For consistency, an array with the same key is
still used in the response when indexing a single file.


A shutdown request may be sent to the indexer from the main server.
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
