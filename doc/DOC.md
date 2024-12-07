## Overview Diagram
![Overview Diagram](https://github.com/DeepAung/gradient/blob/main/doc/Gradient%20Overview%20Diagram.excalidraw.svg)

## TODO

- [x] grader naming ==Gradient==
- [x] Which sql to use (PostgreSQL?)
- [x] Sync this file with Excalidraw
- [ ] Database
  - [ ] Write sql migrate file
  - [ ] Write sql seed file
  - [ ] Setup GCP Bucket
- [ ] Grader server
  - [x] simple gRPC Grade function
  - [ ] result type, memory limit, time limit
  - [ ] actual testcases puller
- [ ] Public server
  - [ ] Users
  - [ ] Auth
  - [ ] Tasks
  - [ ] Submissions

## Future Plans

- [ ] OAuth
- [ ] Favorite the tasks
- [ ] Tasks' tags
- [ ] Text editor
  - [ ] vim-like motion
    - https://github.com/glacambre/firenvim
    - https://github.com/qutebrowser/qutebrowser
  - [ ] code highlighting base on the language
    - https://codemirror.net/
    - https://github.com/highlightjs/highlight.js
- [ ] look at gRPC TLS cert

## Grader Server Usecases

- Grade(code, language) -> stream result
  - _caching testcases in local file storage???_ ==TODO==
  - create directory in name of submission id
  - create code.{language} file
  - run the file with 01.in, 02.in, etc.
  - return output to 01.result
  - check each testcase (01.result and 01.out) (trim result)
  - stream results (P pass, - incorrect, X runtime error, T time limit exceeded) `"PPPPP--XTT"`

## Public Server Usecases

### Users

- GetUser(id)
- UpdateUser(id, updateData)
- UpdatePassword(id, oldPassword, newPassword)
- DeleteUser(id)

### Auth

- SignIn(username, password)
- SignUp(signUpData)
- SignOut(tokenId)
- UpdateTokens(tokenId, refreshToken)

### Tasks

- GetTask(id)
- GetTasks(search, filter, sort, pagination)
- CreateTask(taskData) **admin role**
- UpdateTask(id, taskData) **admin role**
- DeleteTask(id) **admin role**

### Submissions

- SubmitCode(userId, taskId, code, language)
- GetSubmissions(taskId)
- GetSubmissions(userId)
- GetSubmissions(taskId, userId)
- GetSubmission(id)

## Public Server Database

### Users

- id int
- username string
- email string
- password hashString
- picture_url string

### Tokens

- id int
- user_id int
- access_token string
- refresh_token string

### Tasks

- id int
- display_name string
- url_name string
- content_url string
- testcase_count int

### Submissions

- id int
- user_id int
- tasks_id int
- results string `e.g. "PPPPP--XTT" = 5 pass, 2 incorrect, 1 runtime error, 2 time limit exceeded`

## File Storage

- `users/{id}/{filename}` Users' profile picture
- `tasks/{id}/{filename}` Tasks' content
- `testcases/{id}/03.{in|out}` Tasks' Testcases

## Website

### Welcome page

- let ChatGPT generate this for me 😎😎😎
- Can go to **Home page**

### Sign in page

- Sign in and redirect to **Home page**
- Can go to **Sign up page**

### Sign up page

- Sign up then automatically sign in and redirect to **Home page**
- Can go to **Sign in page**

### Home page

- List all tasks
  - Can search, filter(by completed), sort, pagination
  - Can click to **Task detail page**
- Has navbar
  - Can click to **User profile page**
  - Can click to **Admin page** (only **admin role**)

### User profile page

- Show user info
- Update user info. Update password

### Task detail page

- Has problem content pdf available
- Can upload code and submit
  - then get result of each real time from **grader server**
- ==Future==: Text editor with highlights and vim-like-motion
- Can click to view all submissions in this task (to **Submission page**)
- Can click to view my submissions in this task (to **Submission page**)

### Submission page

- Can view submission by task, view my submissions, or both

### Admin page

- list all tasks
- Can go to **Admin task detail page**
- Can go to **Admin task creation page**

### Admin task detail page

- Edit task
- Delete task

### Admin task creation page

- Create task (upload testcases as zip file)
