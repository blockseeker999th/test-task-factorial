# test-task-factorial
Calculates 2 factorial numbers simultaneously

# installation

1. git clone git@github.com:blockseeker999th/test-task-factorial.git (SSH)
   git clone https://github.com/blockseeker999th/test-task-factorial.git (HTTPS)

2. set up applications in docker
``` docker-compose up --build ```
sometimes the application can't wait to load DB-container, in that case do Ctrl+C and then
``` docker-compose up ``` again, we don't need to rebuild the project in that case

3. shutdown application -> Ctrl+C
``` docker-compose down ``` to stop docker-compose containers,networks and volumes

Also you don't need to set up all migrations manually, it'll automatically set up when the docker-compose started

The app calculates the 2 numbers simultaneously with the goroutines help.




