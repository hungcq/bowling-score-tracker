# Bowling Score Tracker
Backend application for tracking bowling scores

## Assumptions
- My interpretation of user story 4:
There is a frame-control mechanism, whereby the current frame can be increased.
Scores of previous frames can't be modified.
- Scores of a player in a frame can be skipped and default to 0 if not entered when changing to next frame.
I think this is more convenient, particularly in the case where a player skips/quits in real life.

## Scope of implementation
- Core game logic, including game rule and score calculation rule
- HTTP endpoints serving features in user stories 1 -> 4
- Error handling for various scenarios (invalid requests, invalid roll results...).
However, due to my limited real-life experience with the game itself,
some scenarios might not be covered or might be handled incorrectly.

### What is not implemented
- Story 5 is not implemented,
given that there is no frontend to show and highlight the current frame and display the final score
- Design the endpoints following REST standard:
the actions allowed in a game is quite specific and limited
(eg increase the current frame, update score of a particular player in the current frame).
Therefore, I just named the endpoints following the action performed.
- Due to limited time to implement and test this system, I didn't implement a storage layer.
Game data is stored transiently in-memory.
However, because the code is already organized by layers and has concerns separated,
a storage layer can easily be implemented and integrated,
by adding a `storage` package that implements a `GameRepository` interface declared in game `managers`.
The `GameManagers` would then load and store the game after each update operation.

## Build & run locally
go build main.go && ./main

## Deployment options
This backend app can be deployed on the cloud as:
### 1. A virtual machine image (eg on AWS)
- Build the VM image
- Create an auto-scaling group
- Set up an application load balancer with HTTPS targeting the auto-scaling group.
Since the app is stateful, route the requests based on hash of IP to make sure 1 session is routed to 1 instance.
- (Optional) set up a domain for the load balancer
### 2. A containerized app (eg AWS ACS)
- Build & push the docker image to registry
- (With a running ECS cluster) Create a task & service definition to run the service
- Config networking for the cluster, including VPC for the cluster, task networking for the service,
then a load balancer targeting the ECS service task. Since the app is stateful,
route the requests based on hash of IP to make sure 1 session is routed to 1 instance.
### 3. A serverless app (eg Lambda function)
- This option requires changing the code to follow the programming model of the service provider.
A storage layer also need to be added to make the app stateless.
- Create & deploy the lambda function
- Create an API gateway to route external requests to the lambda function

## Happy flow & sample request
1. Start a game with player names
```
curl --location 'localhost:80/start_game' \
--header 'Content-Type: application/json' \
--data '{
    "game_type": "TEN_PIN",
    "player_names": [
        "hung",
        "thuy"
    ]
}'

{"game":{"id":2,"current_frame":0,"players":[{"name":"hung","frames":[null,null,null,null,null,null,null,null,null,null],"scores":[0,0,0,0,0,0,0,0,0,0],"total_score":0},{"name":"thuy","frames":[null,null,null,null,null,null,null,null,null,null],"scores":[0,0,0,0,0,0,0,0,0,0],"total_score":0}]}}
```
2. Set score for each player in the current frame of the game
```
curl --location 'localhost:80/1/set_frame_result' \
--header 'Content-Type: application/json' \
--data '{
    "player_index": 1,
    "pins": [
        "X"
    ]
}'

{"game":{"id":1,"current_frame":0,"players":[{"name":"hung","frames":[null,null,null,null,null,null,null,null,null,null],"scores":[0,0,0,0,0,0,0,0,0,0],"total_score":0},{"name":"thuy","frames":[[10],null,null,null,null,null,null,null,null,null],"scores":[10,0,0,0,0,0,0,0,0,0],"total_score":10}]}}

curl --location 'localhost:80/1/set_frame_result' \
--header 'Content-Type: application/json' \
--data '{
    "player_index": 0,
    "pins": [
        "4",
        "4"
    ]
}'

{"game":{"id":1,"current_frame":0,"players":[{"name":"hung","frames":[[4,4],null,null,null,null,null,null,null,null,null],"scores":[8,0,0,0,0,0,0,0,0,0],"total_score":8},{"name":"thuy","frames":[[10],null,null,null,null,null,null,null,null,null],"scores":[10,0,0,0,0,0,0,0,0,0],"total_score":10}]}}
```
3. Increment the current frame of the game
```
curl --location 'localhost:80/1/next_frame' \
--header 'Content-Type: application/json' \

{"game":{"id":1,"current_frame":1,"players":[{"name":"hung","frames":[[4,4],null,null,null,null,null,null,null,null,null],"scores":[8,0,0,0,0,0,0,0,0,0],"total_score":8},{"name":"thuy","frames":[[10],null,null,null,null,null,null,null,null,null],"scores":[10,0,0,0,0,0,0,0,0,0],"total_score":10}]}}
```
4. Repeat step 2 and 3 till the last frame (frame 9)

## Error handling & request sample of all scenarios
See [postman collection](./tracker.postman_collection)

## Code organization
The project follows the port-adapter (hexagonal) architecture, with:
- `core` package: core business logic, with the `managers` interfaces (ports) called by the inbound adapters.
The core consists of domain models with rich behaviors instead of transaction scripts.
- `http_handlers`: inbound adapter for HTTP endpoints
- As mentioned above, there is no a storage layer or any other outbound adapters.
But they could be added by defining outbound ports and implement the adapters.

## Design alternatives
Below is a very brief discussion of the design alternatives and their deployment options.

### Frontend-only application
- The frontend can be anything, from an app in the bowling room to a standard static web.
- If the requirement is simple, this is probably the most attractive option due to is cost-effectiveness and simplicity.
- Deployment option for static web app:
  - Upload to blob storage (eg S3)
  - (Optional) CDN layer on top

### Backend-frontend application
- This option is suitable if we want to evolve the application and add more features.
It also allows storing game/session data.
- Frontend deployment is the same as the front-end option, unless we opt for server-side rendering app
- There are various approaches to deploy the backend, with particular tradeoffs. Here I only list out the main approach:
  - Serverless app
  - Containerized app
  - VMI app
- Storage:
  - Most types of databases can be used, from document to relational to column-wide store
  - The database itself can be serverless or self-managed