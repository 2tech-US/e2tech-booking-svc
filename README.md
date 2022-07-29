# tech2-passenger-svc API document

This is a documentation guide for the Tech2 Microservice Car Sharing system.


## V1

### user

#### Fetch Passenger info:
Authorization token: `ROLE_PASSENGER`, 'ROLE_ADMIN'

```
GET /api/v1/passenger/user/:passengerID
```

### Update Passenger info:
Authorization token: `ROLE_PASSENGER`, 'ROLE_ADMIN'
```
PUT /api/v1/passenger/user/:passengerID
```
Body parameters available (not implement)
