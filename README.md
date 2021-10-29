# postit
## Introduction
Postit is the service that handles posts/comments/likes, it interacts with uacl and notif. Uacl to authorize requests and notif to send notifications to users about messages.

if any request is created and authenticated via the uacl service and the user isn't in the postit db it will create it, since the uacl is the source of truth.

## Production Environment Variables
```
DATABASE_URL - URL for the database. Should also be connecting to the uacl database
VERIFICATION_URL - This is the url of the uacl service, that way it can verify requests
HOST - In case the service needs to run on anything other than 0.0.0.0
PORT - In case the service needs to run on anything other than 80
NOTIFICATION_AUTH - Secret value that notif uses to verify requests
NOTIFICATION_URL - The is the url of the notif service, that way it can send notifications.
EMOTIVES_URL - The is the url of the emotives site, that way it can build notifications.
REDIS_ADDR - The location of Redis so it can connect to it
REDIS_PREFIX - The prefix all keys postit creates have
EMAIL_FROM - Email configuration.
EMAIL_PASSWORD - Email configuration.
EMAIL_LEVEL - What level of logs gets sent to the email address.
ALLOWED_ORIGINS - Cors setup.
```
## Endpoints
```
base URL is postit.emotives.net

all endpoints are user authenticated endpoints, expect for healthz

GET - /healthz - Standard endpoint that just returns ok and a 200 status code. Can be used to test if the service is up
GET - /explore_search - used to fetch posts with query parameters of locations so posts are fetched close within the parameters
POST - /post - create a post based on the request body, see model/post for the request body
GET - /post - Fetch posts, sorted by most recently created, can use page param to find certain pages
GET - /post/{post_id} - Fetch specific post by id
DELETE - /post/{post_id} - Delete specific post by id
POST - /post/{post_id}/like - Creates a like for the specific post by id
DELETE - /post/{post_id}/like/{like_id} - deletes a like for the specific like by id
POST - /post/{post_id}/comment - Creates a comment for the specific post by id. Body request should just be a json object with message key
```
## Database design
Uses a postgres database.
[See here for latest schema, uses the uacl_db](https://github.com/EmotivesProject/databases)