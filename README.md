# Lagoon-sync

Lagoon-sync is part of the Lagoon cli toolset and, indeed, works closely with its parent project.

blah blah ... fill this in.


## Usage

### Syncing a mariadb database

`lagoon-sync sync mariadb --remote-project-name=amazeelabsv4-com --remote-environment-name=dev`

### Syncing a postgres database

`lagoon-sync sync postgres --project-name=drupal9-lando-postgres --source-environment-name=develop`


## Building binary

To build the application and support code into a single binary we need to make sure we configure it so it will run on our images.

Unlike Node.js or Ruby, the Go binary does not need system dependencies such as Go itself in order to run. 

By default, Go will build an executable based on the current system OS and architecture it is ran on (e.g darwin or debian etc). 
Therfore we need to pass this to our `go build` command so it will run in our alpine images:

```
GOOS=linux GOARCH=amd64 go build
```
