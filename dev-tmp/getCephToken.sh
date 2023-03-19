 curl -X POST "https://192.168.79.11:30701/api/auth" \
  -H  "Accept: application/vnd.ceph.api.v1.0+json" \
  -H  "Content-Type: application/json" \
  -d '{"username": "cubeuniverse", "password": "cubeuniverse"}' -k
