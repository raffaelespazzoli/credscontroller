# Spring Legacy Example - A vault unaware application

In this example we are deploying a spring application that is unware of Vault and expect a secret to be available at a given location.

The below picture shows the workflow.

You need to have Vault and Vault-Controller installed as explained [here](../../README.md)

Create a policy that allows the spring-example role to read only from the spring-example generic backend
```
export VAULT_TOKEN=$ROOT_TOKEN
vault policy-write -tls-skip-verify spring-legacy-example ./examples/spring-legacy-example/spring-legacy-example.hcl 
```

Create a secret for the application to consume
```
vault write -tls-skip-verify secret/spring-legacy-example password=pwd 
```

Build the application

```
oc new-project spring-legacy-example
oc new-build registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift~https://github.com/raffaelespazzoli/credscontroller --context-dir=examples/spring-legacy-example --name spring-legacy-example
```
join the network with vault-controller
```
oc adm pod-network join-projects --to vault-controller spring-legacy-example
```
deploy the spring legacy example app
```
oc create -f ./examples/spring-legacy-example/spring-legacy-example.yaml
oc expose svc spring-legacy-example
```
now you should be able to call a service that returns the secret
```
export SPRING_LEGACY_EXAMPLE_ADDR=http://`oc get route | grep -m1 spring | awk '{print $2}'`
curl $SPRING_LEGACY_EXAMPLE_ADDR/secret
```