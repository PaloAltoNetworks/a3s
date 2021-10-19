# a3s

> NOTE: this is a work in progress and this software is not usable yet

a3s (stands for Auth As A Service) is an authentication and authorization server.
It allows to normalize various source of authentications like OIDC, AWS/Azure/GCP Identity tokens into a single and generic authentication tokens that will contain claims (rather than scopes). 
It trusted data like username, groups, age, team, OS, security group or anything really that can been retrieved from the origin auth source.
