# Etsello

Etsello integrates etsy with trello boards. This integration can be used by etsy sellers to link their etsy account with trello. Once linked new purchases will automatically create a card in your trello board in the preferred board and list.

# Running the project
1. .env file contains all the application parameters. Most of them works with the default values except for the following:
*ETSY_CONSUMER_KEY and ETSY_CONSUMER_SECRET*:  Register a new app at [https://www.etsy.com/developers/register](https://www.etsy.com/developers/register). Create an etsy account and setup two factor authentication to register a new app. Use the secret and key for you app as the environment values specified above.
*TRELLO_CONSUMER_KEY and TRELLO_CONSUMER_SECRET*: Onboard a new trello app at [https://trello.com/app-key](https://trello.com/app-key). Use the secret and key as values for the above parameters.
3. Build the project into a docker container
*`docker build -t etsello -f Dockerfile .`*
4. Use docker compose to run the project.
*`docker run --name etsello-mongo-standalone -d -p 27017:27017 -v data:/data/db  mongo:4.2.0-bionic`*
