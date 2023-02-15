# metrics-scope-collector

Google Cloud Monitoringのmetrics scopeに指定したOrganization, Folder以下のProjectを入れるためのもの
https://cloud.google.com/monitoring/settings

## Usage

指定したParent配下の全GCP ProjectのMetrics Scopeを作成する

```
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" {RUN_URL}/metrics-scope-gatherer/create\?parentResourceID\={ID}\&parentResourceType\={organization or folder}
```

Metrics Scopeをすべて削除する

```
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" {RUN_URL}/metrics-scope-gatherer/cleanup
```