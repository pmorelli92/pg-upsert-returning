UPSERT_NOTHING=$(echo '{
  method: "POST",
  url: "http://localhost:8080/upsert-donothing",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "EXECUTING FOR UPSERT DO NOTHING"
jq -ncM "$UPSERT_NOTHING" \
 | vegeta attack -format=json -duration=5s -connections=5 -rate=20 | vegeta encode \
 | vegeta report -type="hist[0,5ms,7ms,9ms,15ms]"


UPSERT_CTE=$(echo '{
  method: "POST",
  url: "http://localhost:8080/upsert-cte",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "EXECUTING FOR UPSERT CTE"
jq -ncM "$UPSERT_CTE" \
 | vegeta attack -format=json -duration=5s -connections=5 -rate=20 | vegeta encode \
 | vegeta report -type="hist[0,5ms,7ms,9ms,15ms]"

UPSERT_CONFLICT=$(echo '{
  method: "POST",
  url: "http://localhost:8080/upsert-conflict",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "EXECUTING FOR UPSERT CONFLICT"
jq -ncM "$UPSERT_CONFLICT" \
  | vegeta attack -format=json -duration=5s -connections=5 -rate=20 | vegeta encode \
  | vegeta report -type="hist[0,5ms,7ms,9ms,15ms]"

UPSERT_LOCK=$(echo '{
  method: "POST",
  url: "http://localhost:8080/upsert-lock",
  body:"{\"id\": \"<ruuid>\"}" | @base64,
  header: {"Content-Type": ["application/json"]}
}' | sed "s/<ruuid>/$(uuidgen | tr '[:upper:]' '[:lower:]')/g")
echo "EXECUTING FOR LOCKS"
jq -ncM "$UPSERT_LOCK" \
  | vegeta attack -format=json -duration=5s -connections=5 -rate=20 | vegeta encode \
  | vegeta report -type="hist[0,5ms,7ms,9ms,15ms]"
