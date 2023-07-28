# ReviveProject

Webserver for Hypixel mods backend. Written with efficiency in mind.

## Why GoLang
Project is created for developers who wish to spend as little money as possible on the hosting server:
- Python and Nodejs are too slow
- Java uses too much memory for JVM
- Rust is fast but too complicated for public to maintain
- Go stands between Java and Rust in the sweet spot

## Download
Obtain latest compiled executable file from GitHub actions. [Click here](https://github.com/aidn3/ReviveProject/actions) to view the latest automated build.  
> Is it safe? The source code is auto complied and built by GitHub servers. You can view the source code on `master` branch. 


## Server Configuration

```
  -cache int
        Time in seconds for global memory cache to speedup repeated requests (default 60)
  -key string
        Hypixel API-Key
  -port int
        Port to listen to (default 8000)
```


## Transformer

A function is provided to modify the final output to customize it for various reasons including:

- Add more information by calculating them on the go
- filter irrelevant information and only show a subset to prevent network overload and prevent abusers from taking
  advantage of the entire output
- simplify the output by processing it in advance allowing clients to have ready to display response

### How to modify
Go to `src/Custom.go` and modify the function:
```go
package src

func Transform(request Request, response *Response) {
    // code here
}
```

Example code:
```go
package src

func Transform(request Request, response *Response) {
  if request.Path == "/skyblock/profiles" {
    response.Data = "{\"success\": true, \"profiles\": \"Modified\"}"
  }
}
```

Json and NBT parsers can also be used for further modification.


## Error Codes

Other status codes are proxied from hypixel.

| Code | Message                                               | Reason                                                                                                                                          |
|------|-------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| 404  | Endpoint not implemented                              | Some default Hypixel endpoints are disabled to prevent abuse                                                                                    |
| 405  | Only GET method is allowed                            | Methods such as POST, PUT, DELETE are disabled                                                                                                  |
| 400  | Parameter(s) required for this endpoint               | Parameter such as `profile`, `uuid`, etc.                                                                                                       |
| 414  | Parameter(s) value is too long                        | If the parameter exceeds 36 character. This usually means the value is malformed. The longest existing Parameter is UUID which is 36 character. |
| 500  | Internal Error encountered. Admins have been notified | An exception appeared that has not been expected                                                                                                |


## Endpoints

- Endpoints such as `/guild` are disabled by default due to it being mostly used by backends bots
- Endpoints such as `/status` and `/skyblock/auction` are disabled by default due to its heavy usage in automated
  tracking
  bots

| Endpoint             | Parameters      | Cache      | Cache Reason                                                                                           |
|----------------------|-----------------|------------|--------------------------------------------------------------------------------------------------------|
| `/player`            | `uuid`          | 12 Hours   | Endpoint mostly used to check overall stats/rank                                                       |
| `/recentgames`       | `uuid`          | 60 Minutes | Endpoint used to view an overview of a player. Stalking isn't required                                 |
| `/skyblock/profiles` | `uuid`          | 5 Minutes  | Used completively for setup checking when partying                                                     |
| `/skyblock/profile`  | `profile`'s id  | 30 Minutes | Longer than `profiles`. Due to its potential abuse with auctions api                                   |
| `/skyblock/museum`   | `profile`'s id  | 12 Hours   | Endpoint used to have an overview of stats                                                             |
| `/skyblock/bingo`    | `uuid`          | 2 Hours    | Can be used to track each other progress. Too fast update will result in abuse with automated trackers |


## Headers

- Extra headers are returned to indicate server health
- No extra headers required to fetch data

| Header      | Explanation                                                                                                                                              |
|-------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `X-age`     | How old the returned data. Timestamp in milliseconds. Useful to indicate the cache age. Current timestamp will be returned if the data is newly fetched. |
| `X-name`    | Server name.                                                                                                                                             |
| `X-version` | Server version.                                                                                                                                          |
