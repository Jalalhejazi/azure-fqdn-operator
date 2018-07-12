# POC Operator: Create AZUR FQDN on service creation

Create a secret with an Azure service principal client credentials, tenant id, and subscription id.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: azure
type: Opaque
data:
  AZURE_CLIENT_ID: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AZURE_CLIENT_SECRET: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AZURE_TENANT_ID: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AZURE_SUBSCRIPTION_ID: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Run the operator with the following manifest.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: azure-fqdn-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: azure-fqdn-operator
  template:
    metadata:
      labels:
        name: azure-fqdn-operator
    spec:
      containers:
        - name: azure-fqdn-operator
          image: neilpeterson/azure-fqdn-operator
          ports:
          - containerPort: 60000
            name: metrics
          command:
          - azure-fqdn-operator
          imagePullPolicy: Always
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: AZURE_CLIENT_ID
              valueFrom:
                secretKeyRef:
                  name: azure
                  key: AZURE_CLIENT_ID
            - name: AZURE_CLIENT_SECRET
              valueFrom:
                secretKeyRef:
                  name: azure
                  key: AZURE_CLIENT_SECRET
            - name: AZURE_TENANT_ID
              valueFrom:
                secretKeyRef:
                  name: azure
                  key: AZURE_TENANT_ID
            - name: AZURE_SUBSCRIPTION_ID
              valueFrom:
                secretKeyRef:
                  name: azure
                  key: AZURE_SUBSCRIPTION_ID
            - name: OPERATOR_NAME
              value: "azure-fqdn-operator"
```

The Azure Go client is then used to update the associated public ip address with an FQDN.

| value | description |
|---|---|
| azure-fqdn-value | DNS prefix name |
| azure-fqdn-rg | AKS node resource group (MC_..) |
| azure-fqdn-location | Resource group location |

Example manifest:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: service-annotation-is-lb-5
  annotations: {
    "azure-fqdn-value": "demo-app-030",
    "azure-fqdn-rg": "MC_scottyCarbone_scottyCarbone_eastus",
    "azure-fqdn-location": "eastus"
  }
spec:
  type: LoadBalancer
  ports:
  - port: 80
```