# [e-inwork.com](https://e-inwork.com)

## Golang Team Indexing Microservice
Golang Team Indexing Microservice is indexing the data of [the Golang Team Microservice](https://github.com/e-inwork-com/go-team-service) in the [Solr](https://solr.apache.org) Search Platform.
To creating a team data, it needs [the Golang User Microservice](https://github.com/e-inwork-com/go-user-service), and basically this microservice needs to run three microservice together in different the Docker container.

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
   docker-compose -f docker-compose.yml up -d
   ```
5. Wait until the status of `curl-local` and `migrate-local` is `exited (0)`, it can check with below command:
   ```
   docker-compose -f docker-compose.yml ps
   ```
6. Create a user in the Golang User Microservice:
   ```
   curl -d '{"email": "jon@doe.com", "password": "pa55word", "first_name": "Jon", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users
   ```
7. Login a user in the Golang User Microservice:
   ```
   curl -d '{"email":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:8000/service/users/authentication
   ```
8. You will get a `token` from the response login, and set it as a token variable for example like below:
   ```
   token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhhY2NkNTUzLWIwZTgtNDYxNC1iOTY0LTA5MTYyODhkMmExOCIsImV4cCI6MTY3MjUyMTQ1M30.S-G5gGvetOrdQTLOw46SmEv-odQZ5cqqA1KtQm0XaL4
   ```
9. Create a team in the Golang Team Microservice for the current user, and change the image to any file image in your folder or just use a sample image in folder `grpc/test/images` :
   ```
   curl -F team_name="Doe's Team" -F team_picture=@/YourRootFolder/.../go-team-indexing-service/grpc/test/images/team.jpg -H "Authorization: Bearer $token"  -X POST http://localhost:8000/service/teams
   ```
10. After successfully creating a team, the Golang Team Microservice will send a team data to the Golang Team Indexing Microservice via gPRC to indexing it in the Solr Search Platform. And the result can check on this link: http://localhost:8983/solr/#/teams/query?q=*:*&q.op=OR&indent=true&useParams=
11. And to run the end to end testing, follow the below commands:
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
12. If you want to debug how this microservice works in editor like [VSCode](https://code.visualstudio.com) for example, just run the Docker Compose for development, before to run `main.go`:
   ```
   docker-compose -f docker-compose.dev.yml up -d
   ```
13. The data of database or indexing will save automatically in the folder `local`, you can delete it if necessary.
14. Have fun!
