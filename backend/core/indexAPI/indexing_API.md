# REST API for indexers

The general standard for the API's requests and responses follow the
[jsend](https://github.com/omniti-labs/jsend) specifications which are very minimal.

HTTP response codes will be ignored, only JSON data will be checked.

`.json` ending in request path as stated by `jsend` will be omitted,
as responses will always be in json-format.

Note that JSON specifications do not support trailing commas.
Adding trailing commas in JSON responses may cause the parser to fail.


Only absolute file paths are used for requests.

A request for indexing a folder path (as opposed to a single file)
should always recursively index sub-folders within the given folder path.


While the main server starts up an indexer, it will check `stderr`.
If the initial ping-request is not able to be responded to,
an error message should be written to `stderr`.


## Ports

Allowed port numbers for indexers are any non-occupied port in the range [39 000, 39 499].

Currently occupied ports are:

```
port 39000, reserved for meta-indexer
port 39001, used by built-in textfile indexer
port 39002, used by built-in PDF indexer
port 39003, used by built-in XML indexer
port 39004, used by built-in web crawler (HTML)
port 39005, used by built-in Golang indexer
```


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
GET /indexfull/FILE_OR_FOLDERPATH
```
Indexer must respond with:
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

**Example:**
Main server requests indexing of a single text file with an active text indexer.
```
GET /indexfull/home/george/my_cool_text_files/mount_vesuvius.txt
```
The text indexer may respond with:
```json
{
    "status": "success",
    "data": {
        "documents":[
            {
                "path": "home/george/my_cool_text_files/mount_vesuvius.txt",
                "source": 0,
                "words": {
                    "volcano": 3,
                    "italy": 1,
                    "the": 65,
                    "erupt": 12,
                } 
            },
        ],
    }
}
```
Note that `"documents"` array only contains a single element.
For consistency, an array with the same key is
still used in the response when indexing a single file.


A shutdown request will be sent to the indexer once all files/folders that
needed to be indexed have been.
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


## Future extentions

In the future, Seekourney Weaver may support handling indexing the difference
between some number of given documents and their new unindexed state.
The request by the main server will be the following:
```
GET /indexdiff/FILEPATH`
```
JSON response follows same format as for `GET /indexfull`
