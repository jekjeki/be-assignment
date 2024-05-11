# Describing about my solution

## Description: 
- I have created user login authentication that need UserId and Password for login it.
- I used third party aurthentication named as JSON Web Token authentication, JSON Web Token is third party that will encrypt our data, in this case example is login data for aunthetication and we can extract it. 
- I have create users can have multiple account such as credit, debit and loan by adding information about account number, and balance for every account. 

I have create executable API with the purpose functions below this: 
 - `/login`, this api use for login authentication of the user with owned credential.
 - `/send`, this api use for send specific balance for another account number with specific user. 
 - `/users/:token`, this api use for get all users token with what token they have, this token getting up fater doing login authentication. After that, the token will be extract it to get specific data that owned. 
 - `/transactions/:userid/:token`, this api used for get all transactions with specific user, if we call this api we can all transaction between withdraw and send data. 
 - `/withdraw/:accountno/:token`, this api use for withdraw the balance of the account. So, if we call the api the account balance will decrease and saved on transactions. 

I have verify every users want to call the api they will needed the token in purpose JWT system to extract it and verify that users is credible. 

I have used postgresSQL database as a database system of my backend. 

I use `docker-compose` for conteneraize the go applications and build it within 2 sub-container as a web and as a database. 
