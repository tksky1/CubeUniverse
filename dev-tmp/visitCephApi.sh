 curl -X POST "https://192.168.79.11:30701/api/" \
  -H  "Accept: application/vnd.ceph.api.v1.0+json" \
  -H  "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJjZXBoLWRhc2hib2FyZCIsImp0aSI6ImVlNTFiNmEwLTU0NzUtNDc0NS1hNzNiLTNiOWJmY2VkYjFiNyIsImV4cCI6MTY3OTA4NTQ0MiwiaWF0IjoxNjc5MDU2NjQyLCJ1c2VybmFtZSI6ImN1YmV1bml2ZXJzZSJ9.HjZP8w2oEbnHgYRVbnTGxsflAl7dTyvabsHjQE6NjHw" \
  -d '{"username": "cubeuniverse", "password": "cubeuniverse"}' -k
