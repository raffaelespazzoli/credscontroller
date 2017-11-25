# Spring Native Example - A Vault aware application using the Vault Kubernetes authentication method

In this example we are deploying a spring application that connects to Vault using the native [kubernetes autentication method](https://www.vaultproject.io/docs/auth/kubernetes.html). Spring vault cloud [supports this method](http://cloud.spring.io/spring-cloud-vault/single/spring-cloud-vault.html#vault.config.authentication.kubernetes) out of the box.
 .

The below picture shows the workflow.

TODO

We assume that vault is alsready installed in the `vault-controller` project. Notice that the vault controller pod does not need to be installed.
If you haven't already done it, configure Vault to use the kubernetes authentication method:
```
oc project vault-controller
oc create sa vault-auth
oc adm policy add-cluster-role-to-user system:auth-delegator vault-auth
secret=`oc describe sa vault-auth | grep 'Tokens:' | awk '{print $2}'`
token=`oc describe secret $secret | grep 'token:' | awk '{print $2}'`
pod=`oc get pods | grep vault | awk '{print $1}'`
oc exec $pod -- cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt >> ca.crt
export VAULT_TOKEN=$ROOT_TOKEN
vault auth-enable -tls-skip-verify kubernetes
vault write -tls-skip-verify auth/kubernetes/config token_reviewer_jwt=$token kubernetes_host=https://kubernetes.default.svc:443 kubernetes_ca_cert=@ca.crt
rm ca.crt
```
notice: `system:auth-delegator` does not seem to work. `cluster-admin` works
Create a policy that allows the spring-example role to read only from the spring-example generic backend
```
export VAULT_TOKEN=$ROOT_TOKEN
vault policy-write -tls-skip-verify spring-example ./examples/spring-native-example/spring-native-example.hcl 
```
Bind the policy with the `spring-native-example` role.
```
vault write -tls-skip-verify auth/kubernetes/role/spring-native-example bound_service_account_names=default bound_service_account_namespaces='*' policies=spring-native-example ttl=1h 
```
with this setup any default service account will be able to use the `spring-native-example` role in vault. 

Create a secret for the application to consume
```
vault write -tls-skip-verify secret/spring-native-example password=pwd 
```

Build the application

```
oc new-project spring-native-example
oc new-build registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift~https://github.com/raffaelespazzoli/credscontroller --context-dir=examples/spring-native-example --name spring-native-example
```
join the network with vault-controller
```
oc adm pod-network join-projects --to vault-controller spring-native-example
```
deploy the spring example app
```
oc create -f ./examples/spring-native-example/spring-native-example.yaml
oc expose svc spring-native-example
```
now you should be able to call a service that returns the secret
```
export SPRING_EXAMPLE_ADDR=http://`oc get route | grep -m1 spring | awk '{print $2}'`
curl $SPRING_EXAMPLE_ADDR/secret
```