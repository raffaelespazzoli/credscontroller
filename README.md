# Vault-based Credential Management Workflow

This example will showcase the workflow illustrated in the following picture

![Vault Controller Flow](../images/vault-controller-raf.png)

1. An Init container requests a wrapped token from the Vault Controller
2. The Vault Controller retrieves the Pod details from the Kubernetes API server
3. If the Pod exists and contains the `vaultproject.io/policies` annotation a unique wrapped token is generated for the Pod.
4. The Vault Controller "callsback" the Pod using the Pod IP obtained from the Kubernetes API.
5. The Init container unwraps the token to obtain a dedicated Vault token.
6. The dedicated token is written to a well-known location and the Init container exits.
7. Another container in the Pod reads the token from the token file.
8. Another container in the Pod renews the token to keep it from expiring.

# Requirements
you need vault CLI installed on your machine

# Create a new project
```
oc new-project vault-controller
```

# Install Vault
```
oc adm policy add-scc-to-user anyuid -z default
oc create configmap vault-config --from-file=vault-config=./vault-config.json
oc create -f vault.yaml
oc create route passthrough vault --port=8200 --service=vault
```
# Initialize Vault
```
export VAULT_ADDR=https://`oc get route | grep -m1 vault | awk '{print $2}'`
vault init -tls-skip-verify -key-shares=1 -key-threshold=1
```
Save the generated key and token. 

# Unseal Vault.
 
You have to repeat this step every time you start vault. 

Don't try to automate this step, this is manual by design. 

You can make the initial seal stronger by increasing the number of keys. 

We will assume that the KEYS environment variable contains the key necessary to unseal the vault and that ROOT_TOKEN contains the root token.

For example:

`export KEYS=M+yDmSrNpFrLuvPYp0q1YTvA+lMaQ6fs0p89i2aKjos=`
`export KEYS=FP3BcfzyE6lIVlMdMKxJGqGVlH+bxCSZO+wwTl1qwiI=`

`export ROOT_TOKEN=8f98666d-f2b3-6756-625e-531744b5101e`
`export ROOT_TOKEN=1d25e02b-0495-2ef4-6344-920f8c024153`

```
vault unseal -tls-skip-verify $KEYS
```

# Install Vault Controller

deploy the vault controller
```
oc create secret generic vault-controller --from-literal vault-token=$ROOT_TOKEN
oc adm policy add-cluster-role-to-user view system:serviceaccount:vault-controller:default
oc create -f vault-controller.yaml
```
# Deploy the Example
```
oc create -f vault-example.yaml
```


# my improvements

1. secure all connections with SSL
2. use an in memory emptyDir to not leave traces of the secret in the node filesystem - done
3. move the authorization labels to a custom API object so that the pod author cannot authorize his pods.
4. support the case where the init-container retrieves the secret as opposed to just a wrapped token that can get the secret (for legacy apps that cannot be modified to talk to Vault
5. add a spring vault example
6. refactor to single command
