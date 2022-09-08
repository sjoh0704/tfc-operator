export CLUSTER_NAME="kubernetes"

# 클러스터 이름을 참조하는 API 서버를 가리킨다.
APISERVER=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$CLUSTER_NAME\")].cluster.server}")

# 토큰 값을 얻는다
TOKEN=$(kubectl get secrets -n tfapplyclaim -o jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='tfc-operator')].data.token}"|base64 --decode)

# TOKEN으로 API 탐색
curl $APISERVER/apis/claim.tmax.io/v1alpha1/namespaces/tfapply/tfapplyclaims/tfapplyclaim-sample/status --header "Authorization: Bearer $TOKEN" --insecure

# API Call을 통한 Status 변경 (status/action: Approve / Plan / Apply)
# Destroy 변경 시 API Call 대신 Manifest 수정 처리 (spec/destroy: true)
curl -X PATCH $APISERVER/apis/claim.tmax.io/v1alpha1/namespaces/tfapply/tfapplyclaims/tfapplyclaim-sample/status --header "Authorization: Bearer $TOKEN" --insecure -H 'Accept: application/json' -H 'Content-Type: application/json-patch+json' -d '[{ "op": "replace", "path": "/status/action", "value": "Approve" }]'
