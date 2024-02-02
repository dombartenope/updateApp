# [Update App Endpoint](https://documentation.onesignal.com/reference/update-an-app)

## The use case for this console application is strictly for updating FCM credentials on apps that are still using the Legacy API. 

### Run `go build && go install` to make this console application accessible from the go/bin globally

### You should run this where you want, or have, a .env file as one is required to store your [User Authentication Key](https://documentation.onesignal.com/docs/keys-and-ids). A local input.json file is also needed in the same directory as where the command is being ran. 

#### There is some error checking enabled, as well as safeguards to confirm that the json file is correct before accepting an app Id to apply this change to. Customers should understand the consequences of doing this are as follows : 
 - Since the customer is uploading a JSON file that does not have a matching sender ID, this will be tied to a different FCM project. This results in new Push Tokens being generated
 - When a new token is created, there is a one hour window where FCM needs to wait before a user can begin receiving push again
