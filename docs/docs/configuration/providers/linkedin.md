---
id: linkedin
title: LinkedIn
---

For LinkedIn, the registration steps are:

1.  Create a new project: https://www.linkedin.com/secure/developer
    - Enter an App name
    - Enter LinkedIn Page
    - Choose App Logo
    - Put checkmark on Legal agreement
    - Click on "Create App"
2.  In the Auth tab:
    - Take note of the **Consumer Key / API Key** and **Consumer Secret / Secret Key**
3.  In the Products tab:
    - Request access for the "Sign In with LinkedIn using OpenID Connect"
    - Wait for the Product to show up in "Added products" section
4.  In the Auth tab
    - Section "OAuth 2.0 scopes" should display scopes "openid", "profile" and "email" at least
