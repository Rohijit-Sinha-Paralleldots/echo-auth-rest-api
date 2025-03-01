Auth REST API

routes

- /auth
    - /signup
        Create a new user
        Input
        - email
        - password
        Output
        - 
    - /signin
        Sign in a user. Successful signin returns an access token and a refresh token
        Input
        - email
        - password
    - /refresh-token
        Get a new access token
    - /signout
        Sign out from a device and invalidate its refresh token
        Input
        - refresh_token
    - /revoke-all-tokens
        Invalidate all refresh tokens
