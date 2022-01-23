Standar Go + Docker Compose + Api in Vault


## Introduction

Hello! My name is Nahuel Ovejero and I'm playing a bit with Go.

While I was investigating Go language I was interested in **clean architecture**. I decided to try implementing **hexagonal architecture (ports and adapters)**.

This API was developing using **Go standard libraries** (uuid is the only exception), but I'm aware that there are tools that will create more efficent and clean code that I'm checking out (like go-resty, gin-gonic, etc).

**I'm pretty happy with the results**: I have learned a **lot**, it was a really fun to play with Go and develop this application. I'm looking forward to continue learning more!


###  Architecture 🔧

![my app diagram](/assets/app-diagram.png)

I named *service* to the left/user side port  and *repository* to the rigth side or app driven ports.

###  Starting up 🚀

```
docker compose up 
```

It will setup all images required, and run the applications tests, and then the it will run the application listening localhost - port 5050:5050.

### Tutorial - Use Examples


### Create

We can create an car (If we don't post the ID it will automatically generate one for us):

```
    curl -i -X POST -H "Content-Type: application/vnd.api+json" \
    http://localhost:5050/cars \
    -d \
    '{
    "data": {
        "type": "cars",
        "id":"ab0bd6f5-c3f5-44b2-b677-acd23cdde73c",
        "owner_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
        "version":1,
    }'
```

**expected response**

```
HTTP/1.1 201 Created - Data of the created Car
HTTP 409 Conflict - Violated Constraints
HTTP 400 Bad Request - Invalid Fields/ Requiered field missing/Invalid uuid
```

### Fetch

We can fetch any car if we got the ID (uuid format):

```
curl -i -X GET http://localhost:5050/cars/ab0bd6f5-c3f5-44b2-b677-acd23cdde73c
```

**expected response**

```
HTTP/1.1 200 OK - Car Data in JSON as response
HTTP 404 Not Found -  Car with given id doesn't exist.
HTTP 400 Bad Request - Invalid uuid
```

### Delete

We can delete an car using the ID. We need to include the owner we want to delete as a query parameter.

```
curl -i -X DELETE http://localhost:5050/cars/ab0bd6f5-c3f5-44b2-b677-acd23cdde73c?version=0
```
**expected response**
```
HTTP/1.1 204 No Content - Sucesfully Deleted
HTTP/1.1 409 Conflict - Incorrect owner
HTTP/1.1 404 Not Found - Car with given uiid was not found
```

---

Nahuel Ovejero
