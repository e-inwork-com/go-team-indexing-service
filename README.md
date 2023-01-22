# [e-inwork.com](https://e-inwork.com)

## Golang Team Indexing Microservice
The Golang Team Indexing Microservice indexes data from the [Golang Team Microservice](https://github.com/e-inwork-com/go-team-service) in the [Solr](https://solr.apache.org) Search Platform. In order to create team data, it relies on the Golang User Microservice. This microservice requires the use of three separate microservices, each running in its own Docker container.

And to run them, follow the steps below:
1. Install Docker
    - https://docs.docker.com/get-docker/
2. Git clone this repository to your folder, and from the terminal run below command:
   ```
   git clone git@github.com:e-inwork-com/go-team-indexing-service
   ```
3. Change the active folder to `go-team-indexing-service`:
   ```
   cd go-team-indexing-service
   ```
4. Run the Docker Compose for the local environment:
   ```
   docker-compose -f docker-compose.local.yml up -d
   ```
5. Wait until the status of `curl-local` and `migrate-local` is `exited (0)` with the below command:
   ```
   docker-compose -f docker-compose.local.yml ps
   ```
6. Create a user in the Golang User Microservice:
   ```
   curl -d '{"email": "jon@doe.com", "password": "pa55word", "first_name": "Jon", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users
   ```
7. Login a user in the Golang User Microservice:
   ```
   curl -d '{"email":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users/authentication
   ```
8. You will get a `token` from the response login, and set it as a token variable for an example like the below:
   ```
   token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhhY2NkNTUzLWIwZTgtNDYxNC1iOTY0LTA5MTYyODhkMmExOCIsImV4cCI6MTY3MjUyMTQ1M30.S-G5gGvetOrdQTLOw46SmEv-odQZ5cqqA1KtQm0XaL4
   ```
9. Create a team in the Golang Team Microservice for the current user. Change the team's image by selecting a file from your local folder or by using a sample image located in `the grpc/test/images` folder. :
   ```
   curl -F team_name="Doe's Team" -F team_picture=@grpc/test/images/team.jpg -H "Authorization: Bearer $token"  -X POST http://localhost:8000/service/teams
   ```
10. After successfully creating a team, the Golang Team Microservice will send the team data to the Golang Team Indexing Microservice via gRPC for indexing in the Solr Search Platform. The indexed results can be viewed by visiting this link: http://localhost:8983/solr/#/teams/query?q=*:*&q.op=OR&indent=true&useParams=
11. And to run the end to end testing, follow the below command:
    ```
    # Down all the current container
    docker-compose -f docker-compose.yml down

    # Run this command if something not working
    docker system prune -a

    # Install all module requirements, install Golang if not installed yet
    go mod tidy

    # Run Docker Compose for test, and wait Wait until the status of "curl-local" and "migrate-local" is "exited (0)"
    docker-compose -f docker-compose.test.yml up -d

    # Run the test
    go test -v ./grpc
    ```
12. To debug how this microservice works in an editor like [VSCode](https://code.visualstudio.com), run the Docker Compose for development before running the `main.go`:
   ```
   docker-compose -f docker-compose.dev.yml up -d
   ```
13. The database and indexing data will be automatically saved in the `local` folder. If necessary, it can be deleted.
14. Have fun!
