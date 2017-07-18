# Scheduler/Calendar (Go)
Basic wireframe for a web app that allows users to schedule time slots. Back-end is written in GoLang and the records are stored in MongoDB. Validations cross check availabilty with records already in the database. Bootstrap for minimal styling.


# Setup
***To Run:***

*Set Environment Variable via Commands or in Bash File*

export PORT="8080"

export MONGODB_URI="mongodb://restofurl"

export MONGODB_DB="mongo_development_database"

  ```
  git clone https://github.com/brianhodges/golang-mongodb-schedule
  cd golang-mongodb-schedule
  ```
  Then run local MongoDB instance
  ```
  mongod
  ```
  
  Finally run app/server
  ```
  go run main.go
  ```
*Then simply navigate in your browser to:* 
 
    http://localhost:8080/


Live Demo: https://golang-mongodb-schedule.herokuapp.com/
