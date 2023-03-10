# hash-sv

## Build Setup 
```bash 
     # build docker image by 
     make build
      
     # run docker container by 
     make run
     
     # Environments:
     HASH_TTL=5m #the time after which the hash will be rebuilt
    
```
## Task

Develop a stateful application containing generated uuid hash in its memory.
Hash should be recreated every 5 minutes.
Application should contain two api servers: gRPC and http.
Each api should implement single endpoint to get actual hash string and hash generation datetime.

Cover code with unit tests where it’s needed.
This app should demonstrate coding quality, app design skills, golang best practices, etc.
It is preferable to upload the results of the work to github or some other public repo.