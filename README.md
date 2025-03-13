 TrainTrack
============


## Functional requirements

- Logging:
  - Training days:
    - Sets: Supersets, Dropsets, Giant sets, Myorep Match sets
    - Reps
    - Weight
    - RPE/RIR
    - Avg. rep time
    - Rest time
  - Training weeks.
  - Training mesocycles.
  - Deloads.
  - Creating templates for the above.
  - Automatic tracking of:
    - Tonnage.
    - Training session duration.
    - Training volume.
- Provide a client/coach framework:
  - A person can send a request to become another person's coach. If the request is accepted, the person is now a coach.
  - One coach - many clients.
  - A client should also be able to be a coach at the same time to some other person.
  - A coach should be able to view/edit a client's training log.
- Subscribe to a person's training log as a watcher. In general there should probably just be a way to choose what rights you grant other people who want to watch your training plan.
- Have a simple chat functionality to facilitate communication between users:
  - Send messages to other users.
  - Attach elements from the log and have them be rendered nicely in the chat.
  - Attach comments on other people's training days, weeks, etc. if u have the rights.
- Analytics:
  - Volume graph.
  - Sets by exercise.
  - Average intensity, while also telling you if you are doing too much.
  - History of performance by exercise
- Other
  - The user should be able to schedule notficiations for himself or for his clients.
  - [PR Scoreboard](https://docs.google.com/spreadsheets/d/1ycoWL810F4lCqcO0AW0wYJgEw_e3GPFtM7GQGbVfevE/edit?gid=657895098#gid=657895098)
  - Search through the exercise gallery.
  - Statistic across all users ranking them by how much muscle they have been gaining.

## Application Logic

### Program Editing Logic

Request data:
Simplest solution is to just send all the data with some updated fields.

1. Client sends a request to `/v1/programs/{program_id}/edit`;
2. The server checks if the client has the rights to edit this program:
  2.1. If not, an Unauthorized error is sent back;
  2.2. If yes, follow through to step 3.
3. The connection is added to the connection hub for this particular program;
4. When the user sends an edit, a lock is acquired on the program, then the edit is applied and sent back to other connections.

The following editing conflicts might arise:
1. User makes an edit to one of the days for example, and saves the edit, but the other user also edited the same block;
2. To resolve this we can prompt the user to decide whether to keep his edits or fetch the updated workout info.

## 🚀 Endpoints

> [!TIP]
> Associate the data being operated upon with the URL path

```sh

# {{{LOGIN/SIGNUP
GET      /v1/login
POST     /v1/signup
# }}}LOGIN/SIGNUP

# {{{ACCESS RIGHTS

GET     /v1/

# }}}ACCESS RIGHTS

# {{{USERS

GET      /v1/users           # get all users
POST     /v1/users           # create new user
GET      /v1/users/{user_id} # get user data by id
DELETE   /v1/users/{user_id} # delete user by id
PATCH    /v1/users/{user_id} # update a user's data

# }}}USERS

# {{{TEMPLATES

GET      /v1/programs                # get all programs from the db NOTE: pass UserId as query param to get his programs only
# POST     /v1/programs                # create new program
GET      /v1/programs/{program_id}   # get program by id
# POST     /v1/programs/{program_id}   # create new program
# PATCH    /v1/programs/{program_id}   # update existing program template
DELETE   /v1/programs/{program_id}   # delete existing program template

GET      /v1/programs/{program_id}/edit   # edit the program through a websocket connection

# }}}TEMPLATES

# {{{LOGS

GET      /v1/logs            # get all logs NOTE: pass UserId as query param to get his programs only
GET      /v1/logs/{log_id}   # get log by id
PATCH    /v1/logs/{log_id}   # update existing log
DELETE   /v1/logs/{log_id}   # delete existing log

GET      /v1/logs/{log_id}/{workout_id}   # get log by id
POST     /v1/logs/{log_id}/{workout_id}   # create new log entry

# }}}LOGS

# {{{CHATS

GET      /v1/chats/            # get all chats
POST     /v1/chats/            # create new chat
GET      /v1/chats/{chat_id}   # get existing chat info
DELETE   /v1/chats/{chat_id}   # delete existing chat

GET      /v1/chats/{chat_id}/messages              # get all chats
POST     /v1/chats/{chat_id}/messages              # create new chat
GET      /v1/chats/{chat_id}/messages/{message_id} # get existing message
PATCH    /v1/chats/{chat_id}/messages/{message_id} # change existing message
DELETE   /v1/chats/{chat_id}/messages/{message_id} # delete existing chat

# }}}CHATS
```

## UI

The exercise UI block should be colored gold if the results beat a previous PR; it should be colored green if it beats the results on the previous microcycle of the training cycle.

- Delay between action and response should be less than 100 milliseconds ideally.

[A delay of less than 100 milliseconds feels instant to a user, but a delay between 100 and 300 milliseconds is perceptible. A delay between 300 and 1,000 milliseconds makes the user feel like a machine is working, but if the delay is above 1,000 milliseconds, your user will likely start to mentally context-switch.](https://designingforperformance.com/performance-is-ux/#:~:text=A%20delay%20of%20less%20than,start%20to%20mentally%20context%2Dswitch.)
## Development requirements

1. Semantic versioning
2. Testing
3. Autogeneratable docs from comments

## References

[Firepad - real time collaborative text editing](https://github.com/FirebaseExtended/firepad)
[Groupware](https://ru.wikipedia.org/wiki/Программное_обеспечение_совместной_работы)
