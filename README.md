# Cookies
## Build and run
### Environment
Create .env file with 2 variables: ADDRESS and LIMIT, where ADDRESS - address on which will be hosted, LIMIT - maximum count of tasks in GET requests result

[VK API for setup bot](https://vk.com/dev/bots_docs "VK's API page")

Example:
```
ADDRESS=127.0.0.1:8080
LIMIT=10
VK_KEY=token
VK_GROUP_ID=group_id
DATABASE=data.db
``` 
### Build
For build use command in project directory:
```
go build .
```
and run:
```
./Cookies
```
## API
### Get tasks (maximum - LIMIT)
For get all tasks, use GET request to ```ADDRESS/api/```

For get tasks from room, use GET request to ```ADDRESS/api/room/{room_id}/```
### Create task
For create new task use POST request to ```ADDRESS/api/add_task/```

Request must include:

```Status: int``` (status of task)


```UserID: int``` (customer ID)

```Room: int``` (customers room number)

## VK API for bots
