/*
Package controller provides HTTP handlers and business logic for managing user authentication, trading pairs, and interactions with various exchanges.

This package is responsible for handling requests related to user accounts, including signing up, logging in, updating passwords, managing trading pairs, retrieving user-found volumes, and deleting user accounts. It interacts with services that manage user data, found volumes, and exchange operations. The controller ensures that all operations are performed securely and efficiently while maintaining the integrity of user data.

Key Features:
  - **User Authentication**: Handles user sign-up and login processes, including password hashing and token generation.
  - **Password Management**: Allows users to update their passwords securely.
  - **Token Management**: Provides functionality to refresh access tokens using refresh tokens.
  - **Account Deletion**: Facilitates the deletion of user accounts along with associated data from exchanges.
  - **Trading Pair Management**: Supports adding, updating, and deleting trading pairs for authenticated users.
  - **Found Volume Retrieval**: Enables retrieval of all found volumes associated with a user's trading pairs.
  - **Error Handling**: Implements robust error handling to provide meaningful feedback to users in case of issues during operations.

Main Components:
  - `userController`: The primary controller that handles requests related to user authentication and trading pairs. It provides methods for signing up users, logging them in, updating passwords, refreshing tokens, managing their trading pairs, and retrieving found volumes.
  - `userPairsController`: Handles requests related to user trading pairs. It provides methods for adding pairs, updating their values, retrieving all user pairs, and deleting specific pairs.

Service Dependencies: The controller relies on several services for its functionality:
  - `UserService`: Manages user-related data and operations.
  - `UserPairsService`: Handles operations related to user trading pairs.
  - `FoundVolumesService`: Manages found volume data associated with trading pairs.
  - `AllExchanges`: Provides access to all exchange instances and their functionalities.

Endpoints:
  - **POST /api/user/auth/signup**: Sign up a new user.
  - **POST /api/user/auth/login**: Authenticate a user and issue tokens if successful.
  - **GET /api/user/auth/tokens**: Retrieve new access and refresh tokens for the authenticated user.
  - **PUT /api/user/auth/password**: Update a user's password.
  - **DELETE /api/user**: Delete the authenticated user's account.
  - **PUT /api/user/pair/update-exact-value**: Update an existing pair for the authenticated user.
  - **POST /api/user/pair**: Add a new trading pair for the authenticated user.
  - **GET /api/user/pair/all-pairs**: Retrieve all pairs for the authenticated user.
  - **GET /api/user/found-volumes**: Retrieve all found volumes associated with the authenticated user's trading pairs.
*/
package controller
