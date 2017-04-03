# Spring Example - A vault aware application

In this example we are deploying a spring application that connects to Vault using the token retrieved by the init container and retrieves a secret.

The below picture shows the workflow.

You need to have Vault and Vault-Controller installed as explained [here](../README.md)

Create a policy that allows the spring-example role to read only from the spring-example generic backend
```
export VAULT_TOKEN=$ROOT_TOKEN
vault policy-write -tls-skip-verify spring-example ./examples/spring-example/spring-example.hcl 
```

Create a secret for the application to consume
```
vault write -tls-skip-verify secret/spring-example password=pwd 
```

Build the application

```
oc new-project spring-example
oc new-build registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift~https://github.com/raffaelespazzoli/credscontroller --context-dir=examples/spring-example --name spring-example
```
join the network with vault-controller
```
oc adm pod-network join-projects --to vault-controller spring-example
```
deploy the spring example app
```
oc create -f ./examples/spring-example/spring-example.yaml
```