# Definitions

## Routes

* Auth
  * GET /auth/callback
    * the callback from a valid authProvider authentication which should return enough info to get a token with identification (code authorization)
  * GET /auth/:authProvider
    * redirects to the authProvider for authentication
* URL
  * Authenticated
    * GET /u
      * Gets all short urls associated with the user
    * POST /u
      * Creates a new short url
    * DELETE /u/:guid
      * Deletes a short url by GUID
  * Unauthenticated
    * GET /u/:guid{8} redirects based on the url stored behind the GUID
      * Normally this would be directly called in the browser to redirect or show a page with the redirect link

## UI
